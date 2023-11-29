import {
    Button,
    FormControl,
    FormLabel,
    Input,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Textarea,
} from '@chakra-ui/react';
import { useContext, useEffect, useRef } from 'react';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import { DesignContext } from '../DesignContext';

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onSave: (data: any) => void;
}

const DesignFormModal = ({ isOpen, onClose, onSave }: Props ) => {
    const initialRef: React.MutableRefObject<null> = useRef(null);
    const { designInEdit } = useContext(DesignContext);

    useEffect(() => {
      const design = designInEdit ? designInEdit : {};
      reset({ ...design })
    }, [designInEdit])

    const schema = yup.object().shape({
      id: yup.string().required(),
      name: yup.string().required(),
      description: yup.string(),
    })

    const { register, handleSubmit, formState: { isValid }, reset } = useForm({
      resolver: yupResolver(schema)
    });

    const handleClose = () => {
      reset();
      onClose();
    }

    return (
    <Modal
      initialFocusRef={initialRef}
      isOpen={isOpen}
      onClose={handleClose}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader textAlign="center">Create Design</ModalHeader>

        <ModalCloseButton />

        <ModalBody pb={6} display="flex" flexDirection="column" gap="10px">
          <FormControl>
            <FormLabel fontSize="12px">ID</FormLabel>
            <Input size="xs" placeholder='ID' {...register('id')} />
          </FormControl>

          <FormControl>
            <FormLabel fontSize="12px">Name</FormLabel>
            <Input size="xs" placeholder='Name' {...register('name')}/>
          </FormControl>

          <FormControl mt={4}>
            <FormLabel fontSize="12px">Description</FormLabel>
            <Textarea size="xs" placeholder='Description' {...register('description')}/>
          </FormControl>
        </ModalBody>

        <ModalFooter>
          <Button onClick={handleSubmit(onSave)} colorScheme='blue' mr={3} isDisabled={!isValid}>
            Save
          </Button>
          <Button onClick={handleClose}>Cancel</Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default DesignFormModal;