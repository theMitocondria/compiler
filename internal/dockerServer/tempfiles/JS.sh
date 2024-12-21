#Tests Done :
    # AC
    # Compile Eror done
    # TLE

#Security checks remaining

cat << EOF > temp.js

EOF


ulimit -v 254800 

OUTPUT=$(echo "1000000" | timeout 1s node check.js 2>&1)
EXIT_CODE=$?  # Capture the exit code of the last command

if [ $EXIT_CODE -eq 143  ]; then
    OUTPUT="TLE"
elif [ $EXIT_CODE -eq 139 ]; then
    OUTPUT="Error: Segmentation fault (invalid memory access)."
elif [ $EXIT_CODE -eq 136 ]; then
    OUTPUT="Error: Floating-point exception."
elif [ $EXIT_CODE -eq 134 ]; then
    OUTPUT="Error: Program aborted."
elif [ $EXIT_CODE -eq 127 ]; then
    OUTPUT="Error: Executable not found."
elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
    SIGNAL=$((EXIT_CODE - 128))
    OUTPUT="Error: Program terminated by signal $SIGNAL."
fi

if [ -f temp.js ];then 
    rm temp.js
fi


echo "$OUTPUT , $EXIT_CODE"
