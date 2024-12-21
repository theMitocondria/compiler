#Tests done : 
    # TLE ERROR
    # No Error code (AC)
    # Syntax Error 
    # Runtime error (divide by zero )
    # Unlimited stack calls (stack space)
    # system commands disabled
    # eval disabled, exec disabled
    # get requests banned (implement in C++ also) (tested but not sure)

# Crash Scenario: 
    # if MLE done (then container has crashed, restart it immediately)

cat << 'EOF' > temp.py
import builtins
import os

builtins.eval = None
builtins.exec = None
builtins.open = None
os.setuid = None 
os.system = None
os.stetgid = None

# User code goes here
print(100)
EOF

# Syntax check the Python code
python3 -m py_compile temp.py 2>/tmp/syntax_error.log
EXIT_CODE=$?  # Capture the exit code of the syntax check

if [ $EXIT_CODE -ne 0 ]; then
    # If there's a syntax error, display it
    cat /tmp/syntax_error.log
else
    # Set a memory limit and run the Python script with a timeout
    ulimit -v 284800  # Limit memory to 284 MB
    OUTPUT=$(echo "25" | timeout 1s python3 temp.py 2>&1)
    EXIT_CODE=$?

    # Handle different exit codes
    if [ $EXIT_CODE -eq 143 ]; then
        OUTPUT="TLE"
    elif [ $EXIT_CODE -eq 1 ]; then
        OUTPUT="Error: Runtime Error"
    elif [ $EXIT_CODE -eq 139 ]; then
        OUTPUT="Error: Memory Limit Exceeded"
    elif [ $EXIT_CODE -ne 0 ]; then
        OUTPUT="Error: Program terminated unexpectedly. $EXIT_CODE"
    else
        TLE="false"
    fi

    # Cleanup the Python script
    rm temp.py
    # Output the result
<<<<<<< HEAD:FinalLang/PYTHON/PYTHON.sh
    echo "$OUTPUT , $EXIT_CODE"
fi


=======
    echo "$OUTPUT"
fi
>>>>>>> b9d3caae3d8bafae8ed17eaf27e9332c194fc013:internal/dockerServer/PYTHON/PYTHON.sh
