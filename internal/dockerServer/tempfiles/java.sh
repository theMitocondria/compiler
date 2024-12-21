# Tests done : 
#   AC code
#   Compilation Error
#   TLE 

# Special case in Java: 
#   no segmentation fault is there in Java
#   no runtime error
#   no safety measures taken

# Read user code from input
USER_CODE="$1"

# Extract the class name containing the main method
CLASS_NAME=$(echo "$USER_CODE" | grep -oP 'public\s+class\s+\K\w+(?=.*\bpublic\s+static\s+void\s+main\b)')

if [ -z "$CLASS_NAME" ]; then
    echo "Error: No class with a main method found."
    exit 1
fi

# Create the Java file with the correct name
echo "$USER_CODE" > "$CLASS_NAME.java"
echo "$USER_INPUT" > input.txt
# Compile the Java file
javac "$CLASS_NAME.java" 2> compile_errors.txt
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    ERROR=$(cat compile_errors.txt)
else 
    echo 'Main-Class: '"$CLASS_NAME" > manifest.txt

    jar cfm "$CLASS_NAME.jar" manifest.txt "$CLASS_NAME.class"

    OUTPUT=$(timeout 1s /usr/lib/jvm/java-17-openjdk/bin/java -XX:+UseSerialGC -XX:TieredStopAtLevel=1 -XX:NewRatio=5 -Xms8M -Xmx128M -Xss64M -DONLINE_JUDGE=true -jar "$CLASS_NAME.jar" < input.txt 2>&1)
    EXIT_CODE=$?  # Capture the exit code of the last command

    if [ $EXIT_CODE -eq 124 ]; then
        ERROR="TLE"
        OUTPUT=""
    elif [ $EXIT_CODE -eq 136 ]; then
        ERROR="Error: Floating-point exception."
        OUTPUT=""
    elif [ $EXIT_CODE -eq 134 ]; then
        ERROR="Error: Program aborted."
        OUTPUT=""
    elif [ $EXIT_CODE -eq 126 ]; then
        ERROR="Error: Permission denied (unable to execute)."
        OUTPUT=""
    elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
        SIGNAL=$((EXIT_CODE - 128))
        ERROR="Error: Program terminated by signal $SIGNAL."
        OUTPUT=""
    elif [ $EXIT_CODE -ne 0 ]; then
        ERROR="Error: Exit code $EXIT_CODE."
        OUTPUT=""
    fi

    # Clean up temporary files
    rm -f "$CLASS_NAME.java" "$CLASS_NAME.class" "$CLASS_NAME.jar" manifest.txt compile_errors.txt input.txt

    echo "OUTPUT : $OUTPUT \n ERROR : $ERROR"
fi






