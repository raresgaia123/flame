/**
 * Copyright 2023 Cisco Systems, Inc. and its affiliates
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import { useQuery } from "@tanstack/react-query";
import { LOGGEDIN_USER } from "../../../constants";
import { Job } from "../../../entities/Job";
import ApiClient from "../../../services/api-client";

const useJobStatus = (id: string | undefined) => {
    const apiClient = new ApiClient<Job>(`users/${LOGGEDIN_USER.name}/jobs/${id}/status`);
    return useQuery({
        enabled: !!id,
        queryKey: ['jobStatus', id],
        queryFn: apiClient.getAll,
    });
}

export default useJobStatus;