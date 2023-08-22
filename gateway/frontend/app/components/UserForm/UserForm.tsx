'use client'

/* Instruments */
import {
   useSelector,
   useDispatch,
   setName,
   setEmail,
   selectName,
   selectEmail,
} from '@/lib/redux'

import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';


export const UserForm = () => {
  const dispatch = useDispatch()

  const name = useSelector(selectName)
  const email = useSelector(selectEmail)

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch(setName(e.target.value))
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch(setEmail(e.target.value))
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Form Submitted with:", { name, email })
  };

  return (
    <form onSubmit={handleSubmit}>
      <TextField
        label="Name"
        variant="outlined"
        value={name}
        onChange={handleNameChange}
      />
      <TextField
        label="Email"
        type="email"
        variant="outlined"
        value={email}
        onChange={handleEmailChange}
      />
      <Button variant="contained" color="primary" type="submit">
        Submit
      </Button>
    </form>
  )
}
