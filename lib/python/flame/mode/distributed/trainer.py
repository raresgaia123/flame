# Copyright 2022 Cisco Systems, Inc. and its affiliates
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0
"""distributed FL trainer."""

# import hashlib
import logging
from collections import OrderedDict

from ...channel_manager import ChannelManager
from ...common.custom_abcmeta import ABCMeta, abstract_attribute
from ...common.util import (MLFramework, get_ml_framework_in_use,
                            mlflow_runname, valid_frameworks)
from ...registries import registry_provider
from ..composer import Composer
from ..message import MessageType
from ..role import Role
from ..tasklet import Loop, Tasklet

logger = logging.getLogger(__name__)

TAG_RING_ALLREDUCE = 'ring_allreduce'

class Trainer(Role, metaclass=ABCMeta):
    """Trainer implements an ML training role."""

    @abstract_attribute
    def config(self):
        """Abstract attribute for config object."""

    @abstract_attribute
    def dataset_size(self):
        """Abstract attribute for size of dataset used to train."""

    def internal_init(self) -> None:
        """Initialize internal state for role."""
        self.cm = ChannelManager()
        self.cm(self.config)
        self.cm.join_all()

        self.registry_client = registry_provider.get(self.config.registry.sort)
        # initialize registry client
        self.registry_client(self.config.registry.uri, self.config.job.job_id)

        base_model = self.config.base_model
        if base_model and base_model.name != "" and base_model.version > 0:
            self.model = self.registry_client.load_model(
                base_model.name, base_model.version)
        self.ring_weights = None # store the latest model weights from ring all-reduce

        self.registry_client.setup_run(mlflow_runname(self.config))
        self.metrics = dict()

        self._round = 1
        self._rounds = self.config.hyperparameters['rounds']
        self._work_done = False

        self.is_committer = False

        self.framework = get_ml_framework_in_use()
        if self.framework == MLFramework.UNKNOWN:
            raise NotImplementedError(
                "supported ml framework not found; "
                f"supported frameworks are: {valid_frameworks}")

        if self.framework == MLFramework.PYTORCH:
            self._scale_down_weights_fn = self._scale_down_weights_pytorch
            self._get_send_chunk_fn = self._get_send_chunk_pytorch
            self._allreduce_fn = self._allreduce_pytorch
            self._allgather_fn = self._allgather_pytorch

        elif self.framework == MLFramework.TENSORFLOW:
            self._scale_down_weights_fn = self._scale_down_weights_tensorflow
            self._get_send_chunk_fn = self._get_send_chunk_tensorflow
            self._allreduce_fn = self._allreduce_tensorflow
            self._allgather_fn = self._allgather_tensorflow

    def _ring_allreduce(self, tag: str) -> None:
        if tag != TAG_RING_ALLREDUCE:
            return

        channel = self.cm.get_by_tag(tag)
        if not channel:
            logger.debug(f"channel not found with tag {tag}")
            return

        success, total_data_count = self._member_check(channel)
        logger.debug(f"member check: {success}")
        logger.debug(f"total_data_count: {total_data_count}")
        if not success:
            # members don't agree, we can't do ring-allreduce
            return

        self._update_weights()

        self._scale_down_weights_fn(total_data_count)

        self._do_ring_allreduce(channel)

        self._update_model()

    def _do_ring_allreduce(self, channel):
        # This method is implemented based on
        # https://github.com/baidu-research/baidu-allreduce/blob/master/collectives.cu
        logger.info("starting ring-allreduce")
        my_id = channel.get_backend_id()
        ends = channel.ends()
        ends.append(my_id)
        ends.sort()

        rank = ends.index(my_id)
        size = len(ends)

        logger.debug(f"weights length = {len(self.weights)}")
        chunk_size = int(len(self.weights) / size)
        chunk_sizes = [chunk_size] * size

        residual = len(self.weights) % size
        for i in range(residual):
            chunk_sizes[i] += 1

        chunk_ends = [0] * size
        chunk_ends[0] = chunk_sizes[0]
        for i in range(1, len(chunk_ends)):
            chunk_ends[i] = chunk_sizes[i] + chunk_ends[i - 1]

        recv_from = (rank - 1 + size) % size
        send_to = (rank + 1) % size

        # enable the following "digest" computation lines
        # to check if ring-allreduce work correctly.
        #
        # digest = hashlib.sha1(str(self.weights).encode('utf-8')).hexdigest()
        # logger.debug(f"initial: weight digest - {digest}")

        # allreduce
        for i in range(size - 1):
            send_chunk_idx = (rank - i + size) % size
            from_idx = chunk_ends[send_chunk_idx] - chunk_sizes[send_chunk_idx]
            to_idx = chunk_ends[send_chunk_idx]

            send_chunk = self._get_send_chunk_fn(from_idx, to_idx)
            logger.debug(
                f"sending chunk: {len(send_chunk)} to {ends[send_to]}")

            channel.send(ends[send_to], {MessageType.WEIGHTS: send_chunk})

            recv_chunk_idx = (rank - i - 1 + size) % size
            msg = channel.recv(ends[recv_from])
            while MessageType.WEIGHTS not in msg:
                msg = channel.recv(ends[recv_from])
            recv_chunk = msg[MessageType.WEIGHTS]

            logger.debug(
                f"receiving chunk: {len(recv_chunk)} from {ends[recv_from]}")

            try:
                assert len(recv_chunk) == chunk_sizes[recv_chunk_idx]
            except AssertionError:
                logger.error(
                    f"AssertionError: got {recv_chunk} from {ends[recv_from]}")
                exit(1)

            from_idx = chunk_ends[recv_chunk_idx] - chunk_sizes[recv_chunk_idx]

            self._allreduce_fn(from_idx, recv_chunk)

        # digest = hashlib.sha1(str(self.weights).encode('utf-8')).hexdigest()
        # logger.debug(f"after allreduce: weight digest - {digest}")

        # allgather
        for i in range(size - 1):
            send_chunk_idx = (rank - i + 1 + size) % size
            from_idx = chunk_ends[send_chunk_idx] - chunk_sizes[send_chunk_idx]
            to_idx = chunk_ends[send_chunk_idx]

            send_chunk = self._get_send_chunk_fn(from_idx, to_idx)
            channel.send(ends[send_to], {MessageType.WEIGHTS: send_chunk})

            recv_chunk_idx = (rank - i + size) % size
            msg = channel.recv(ends[recv_from])
            recv_chunk = msg[MessageType.WEIGHTS]

            try:
                assert len(recv_chunk) == chunk_sizes[recv_chunk_idx]
            except AssertionError:
                logger.error(
                    f"AssertionError: got {recv_chunk} from {ends[recv_from]}")
                exit(1)

            from_idx = chunk_ends[recv_chunk_idx] - chunk_sizes[recv_chunk_idx]

            self._allgather_fn(from_idx, recv_chunk)

        # digest = hashlib.sha1(str(self.weights).encode('utf-8')).hexdigest()
        # logger.debug(f"after allgather: weight digest - {digest}")
        logger.info("finished ring-allreduce")

    """BEGIN: pytorch functions"""

    def _scale_down_weights_pytorch(self, total: int) -> None:
        if total == 0:
            return

        rate = self.dataset_size / float(total)

        self.weights = {k: v * rate for k, v in self.weights.items()}

    def _get_send_chunk_pytorch(self, from_idx, to_idx):
        send_chunk = OrderedDict()
        for i, key in enumerate(self.weights):
            if i < from_idx:
                continue
            elif i >= to_idx:
                break

            send_chunk[key] = self.weights[key]

        return send_chunk

    def _allreduce_pytorch(self, from_idx, recv_chunk):
        # from_idx is not used in case of pytorch
        # recv_chunk is an ordered dictionary
        # so, it's okay to update weights based on keys in recv_chunk
        for k, v in recv_chunk.items():
            self.weights[k] += v

    def _allgather_pytorch(self, from_idx, recv_chunk):
        # from_idx is not used in case of pytorch
        # recv_chunk is an ordered dictionary
        # so, it's okay to update weights based on keys in recv_chunk
        for k, v in recv_chunk.items():
            self.weights[k] = v

    """END: pytorch functions"""
    """BEGIN: tensorflow functions"""

    def _scale_down_weights_tensorflow(self, total: int) -> None:
        if total == 0:
            return

        rate = self.dataset_size / float(total)
        self.weights = [weight * rate for weight in self.weights]

    def _get_send_chunk_tensorflow(self, from_idx, to_idx):
        return self.weights[from_idx:to_idx]

    def _allreduce_tensorflow(self, from_idx, recv_chunk):
        for i in range(len(recv_chunk)):
            self.weights[from_idx + i] += recv_chunk[i]

    def _allgather_tensorflow(self, from_idx, recv_chunk):
        for i in range(len(recv_chunk)):
            self.weights[from_idx + i] = recv_chunk[i]

    """END: tensorflow functions"""

    def _handle_member_check(self, channel, end, digest) -> tuple[bool, int]:
        """Handle member check message.

        Returns
        -------
        success: boolean variable to tell if there is a consistent group or not
        dataset_size: the size of dataset used by a member
        """
        dataset_size = 0
        while True:
            msg = channel.peek(end)
            if msg is not None and MessageType.WEIGHTS in msg:
                logger.debug("weights msg seen; let's stop member check")
                break

            msg = channel.recv(end)

            # check if a new trainer needs the latest weights
            if MessageType.NEW_TRAINER in msg and self.is_committer:
                logger.debug(f"{channel.get_backend_id()} sending weights to the new trainer {end}")
                channel.send(end, {MessageType.RING_WEIGHTS: self.ring_weights})

            logger.debug(f"end_id: {end}, msg: {msg}")
            if MessageType.MEMBER_DIGEST not in msg:
                logger.debug("no member digest found")
                return False, 0
            if MessageType.DATASET_SIZE not in msg:
                logger.debug("no dataset size found")
                return False, 0
            if MessageType.ROUND not in msg:
                logger.debug("no round info found")
                return False, 0

            self._round = max(self._round, msg[MessageType.ROUND])

            other_digest = msg[MessageType.MEMBER_DIGEST]
            if digest != other_digest:
                logger.debug(f"mine: {digest}, other: {other_digest}")
                return False, 0

            dataset_size = msg[MessageType.DATASET_SIZE]

            # new trainer fetches the latest weights from the committer
            if msg[MessageType.IS_COMMITTER]:
                while self.ring_weights is None:
                    msg = channel.recv(end)
                    logger.debug(f"new trainer {channel.get_backend_id()} fetching weights from {end} ")
                    if MessageType.RING_WEIGHTS in msg:
                        self.weights = msg[MessageType.RING_WEIGHTS]
                        self._update_model()
                        break

            if channel.is_rxq_empty(end):
                break

        return True, dataset_size

    def _member_check(self, channel) -> tuple[bool, int]:
        digest = channel.ends_digest()
        if digest == "":
            # This means that there is no ends in the channel,
            # so there is no point of broadcasting digest.
            # If the empty digest is broadcast, it can cause a bug
            self.is_committer = True
            self._update_weights()
            self._update_model()
            return False, 0

        msg = {
            MessageType.MEMBER_DIGEST: digest,
            MessageType.DATASET_SIZE: self.dataset_size,
            MessageType.ROUND: self._round,
            MessageType.IS_COMMITTER: self.is_committer
        }
        if self.ring_weights is None:
            logger.debug("Sending arrival message...")
            msg[MessageType.NEW_TRAINER] = True
        logger.debug(f"member check msg = {msg}")
        channel.broadcast(msg)

        total_count = self.dataset_size
        ends = channel.ends()
        for end in ends:
            success, size = self._handle_member_check(channel, end, digest)
            if not success:
                logger.debug(f"_handle_member_check failed for {end}")
                return False, 0

            total_count += size

        my_taskid = channel.get_backend_id()
        ends.append(my_taskid)
        ends.sort()
        # if my taskid is in the first ends, then it's selected as a committer
        self.is_committer = True if my_taskid in ends[:1] else False

        # check if work is done by others
        # if work is done, then no further distributed learning needed
        self._work_done = (self._round > self._rounds)
        if self._work_done:
            return False, 0

        return True, total_count

    def get(self, tag: str) -> None:
        """Get data from remote role(s)."""
        return

    def put(self, tag: str) -> None:
        """Set data to remote role(s)."""
        return

    def save_metrics(self):
        """Save metrics in a model registry."""
        logger.debug(f"saving metrics: {self.metrics}")
        if self.metrics:
            self.registry_client.save_metrics(self._round, self.metrics)
            logger.debug("saving metrics done")

    def update_metrics(self, metrics: dict[str, float]):
        """Update metrics."""
        self.metrics = self.metrics | metrics

    def _update_model(self):
        """Update model with weights."""
        if self.framework == MLFramework.PYTORCH:
            self.model.load_state_dict(self.weights)
        elif self.framework == MLFramework.TENSORFLOW:
            self.model.set_weights(self.weights)
        self.ring_weights = self.weights

    def _update_weights(self):
        """Save weights from model."""
        if self.framework == MLFramework.PYTORCH:
            self.weights = self.model.state_dict()
        elif self.framework == MLFramework.TENSORFLOW:
            self.weights = self.model.get_weights()

    def increment_round(self):
        """Increment the round counter."""
        logger.debug(f"Incrementing current round: {self._round}")

        self._round += 1
        self._work_done = (self._round > self._rounds)

        channel = self.cm.get_by_tag(TAG_RING_ALLREDUCE)
        if not channel:
            logger.debug(f"channel not found for tag {TAG_RING_ALLREDUCE}")
            return

        # set necessary properties to help channel decide how to select ends
        channel.set_property("round", self._round)

    def save_params(self):
        """Save hyperparamets in a model registry."""
        logger.debug(f"saving params: is_commiter: {self.is_committer}")
        if self.config.hyperparameters and self.is_committer:
            self.registry_client.save_params(self.config.hyperparameters)

    def save_model(self):
        """Save model in a model registry."""
        logger.debug(f"saving model: is_commiter: {self.is_committer}")
        if self.model and self.is_committer:
            model_name = f"{self.config.job.name}-{self.config.job.job_id}"
            self.registry_client.save_model(model_name, self.model)

    def compose(self) -> None:
        """Compose role with tasklets."""
        with Composer() as composer:
            self.composer = composer

            task_internal_init = Tasklet(self.internal_init)

            task_load_data = Tasklet(self.load_data)

            task_init = Tasklet(self.initialize)

            task_allreduce = Tasklet(self._ring_allreduce, TAG_RING_ALLREDUCE)

            task_train = Tasklet(self.train)

            task_eval = Tasklet(self.evaluate)

            task_increment_round = Tasklet(self.increment_round)

            task_save_metrics = Tasklet(self.save_metrics)

            task_save_params = Tasklet(self.save_params)

            task_save_model = Tasklet(self.save_model)

            # create a loop object with loop exit condition function
            loop = Loop(loop_check_fn=lambda: self._work_done)
            task_internal_init >> task_load_data >> task_init >> loop(
                task_train >> task_allreduce >> task_eval >> task_save_metrics
                >> task_increment_round) >> task_save_params >> task_save_model

    def run(self) -> None:
        """Run role."""
        self.composer.run()

    @classmethod
    def get_func_tags(cls) -> list[str]:
        """Return a list of function tags defined in the trainer role."""
        return [TAG_RING_ALLREDUCE]
