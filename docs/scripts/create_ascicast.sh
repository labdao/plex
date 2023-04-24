#!/bin/bash
# run inside container image with asciinema installed

# File containing the list of CLI commands to record
commands_file="scripts/ascicast-commands.txt"

# Display the commands in the file
while IFS= read -r cmd; do
  echo "Command: $cmd"
done < "$commands_file"

# Read commands from the file into an array
commands=()
while IFS= read -r cmd; do
  commands+=("$cmd")
done < "$commands_file"

# Loop through the commands array and create a recording for each one
for cmd in "${commands[@]}"; do
  # Create a unique filename for each recording
  filename="$(echo "$cmd" | tr -d '[:space:]/' | tr -cd '[:alnum:]._-').cast"

  # Create a script to simulate typing the command character by character
  typed_cmd_script="tmp.sh"
  echo "#!/bin/bash" > "$typed_cmd_script"
  for ((i=0; i<${#cmd}; i++)); do
    echo "printf '%s' '${cmd:$i:1}'" >> "$typed_cmd_script"
    echo "sleep 0.1" >> "$typed_cmd_script"
  done
  echo "printf '\n'" >> "$typed_cmd_script"
  echo "$cmd" >> "$typed_cmd_script"
  echo "exit" >> "$typed_cmd_script"
  chmod +x "$typed_cmd_script"

  # Start the recording, execute the command, and then exit the shell
  asciinema rec -c "bash $typed_cmd_script" $filename -y -i 2 --overwrite 

  # Cleanup the temporary script
  rm -f "$typed_cmd_script"
done