
#stillexisting problems : 
#   system command execute ho jari h 
#   file system access ho paara h
#   block get and post rquests
#   use of chroot
#   exit testing 
#   MLE
#   Output End
#tests done :
#   Simple code no error
#   Compilation Error show krna
#   TLE Encounter
#   Segmentation Encounter
#   Runtime / Floating point  

TEMP_CPP=$(mktemp /tmp/temp.XXXXXX.cpp)
TEMP_EXE=$(mktemp /tmp/temp.XXXXXX)
TEMP_INPUT=$(mktemp /tmp/temp.XXXXXX.input)

cat << EOF > $TEMP_CPP
//some code by user .....
EOF

cat << EOF > $TEMP_INPUT
here comes the user input
EOF

ERROR=$(g++ -DONLINE_JUDGE=true -O2 -Wall -Wextra -Werror -std=c++17 -o $TEMP_EXE $TEMP_CPP 2>&1)



if [ ! -f $TEMP_EXE ]; then
    echo "Compilation failed with the following error:"
else 
    ulimit -v 254800 
    ulimit -t 5
    ulimit -f 1000

    OUTPUT=$(timeout 1s $TEMP_EXE < $TEMP_INPUT 2>&1)
    EXIT_CODE=$?  # Capture the exit code of the last command

    if [ $EXIT_CODE -eq 143  ]; then
        OUTPUT=""
        ERROR="Error: Time limit exceeded."
    elif [ $EXIT_CODE -eq 139 ]; then
        OUTPUT=""
        ERROR="Error: Segmentation fault."
    elif [ $EXIT_CODE -eq 136 ]; then
        OUTPUT=""
        ERROR="Error: Floating point exception."
    elif [ $EXIT_CODE -eq 134 ]; then
        OUTPUT=""
        ERROR="Error: Aborted."
    elif [ $EXIT_CODE -eq 127 ]; then
        OUTPUT=""
        ERROR="Error: Command not found."
    elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
        SIGNAL=$((EXIT_CODE - 128))
        OUTPUT=""
        ERROR="Error: Interrupted with signal $SIGNAL."
    elif [ $EXIT_CODE -ne 0 ]; then
        OUTPUT=""
        ERROR="Error: Exit code $EXIT_CODE."
    fi

fi

rm -f $TEMP_CPP $TEMP_EXE $TEMP_INPUT

echo "$OUTPUT "
echo "$ERROR"