#Tests done : 
    # AC code
    # Compilation Error
    # TLE 
    
#special case in Java: 
    # no segmentation fault is there in Java
    # no runtime error
    # no safety measures taken

cat << EOF > Main.java
public class Main {
    public static void main(String[] args) {
        // This will cause an OutOfMemoryError
        int a = 15/3;
        System.out.print(a);
    }
}

EOF

javac Main.java 2> compile_errors.txt
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    echo "Compilation Error:"
    cat compile_errors.txt
else 
    echo 'Main-Class: Main'>manifest.txt

    jar cfm Main.jar manifest.txt Main.class

    OUTPUT=$(echo "10000000" | timeout 1s /usr/lib/jvm/java-17-openjdk/bin/java -XX:+UseSerialGC -XX:TieredStopAtLevel=1 -XX:NewRatio=5 -Xms8M -Xmx128M -Xss64M -DONLINE_JUDGE=true -jar Main.jar  2>&1)
    EXIT_CODE=$?  # Capture the exit code of the last command

    if [ $EXIT_CODE -eq 143  ]; then
        OUTPUT="TLE"
    elif [ $EXIT_CODE -eq 136 ]; then
        OUTPUT="Error: Floating-point exception."
    elif [ $EXIT_CODE -eq 134 ]; then
        OUTPUT="Error: Program aborted."
    elif [ $EXIT_CODE -eq 126 ]; then
        OUTPUT="Error: Permission denied (unable to execute)."
    elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
        SIGNAL=$((EXIT_CODE - 128))
        OUTPUT="Error: Program terminated by signal $SIGNAL."
    elif [ $EXIT_CODE -ne 0 ]; then
        OUTPUT="Error: Program terminated unexpectedly."
    fi

    if [ -f Main.class ]; then
        rm Main.class
    fi

    if [ -f Main.jar ]; then
        rm Main.jar
    fi

    if [ -f Main.java ]; then
        rm Main.java
    fi

    if [ -f manifest.txt ]; then
        rm manifest.txt
    fi

    echo "$OUTPUT , $EXIT_CODE"

fi