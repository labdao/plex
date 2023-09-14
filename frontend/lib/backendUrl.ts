export default function backendUrl() {
  if(!process.env.NEXT_PUBLIC_BACKEND_URL) {
    throw new Error("The environment variable NEXT_PUBLIC_BACKEND_URL must be defined.");
  }
  return process.env.NEXT_PUBLIC_BACKEND_URL
}
