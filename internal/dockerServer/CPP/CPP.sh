# #stillexisting problems : 
# #   system command execute ho jari h 
# #   file system access ho paara h
# #   block get and post rquests
# #   use of chroot
# #   exit testing 
# #   MLE
# #   Output End
# #tests done :
# #   Simple code no error
# #   Compilation Error show krna
# #   TLE Encounter
# #   Segmentation Encounter
# #   Runtime / Floating point  

# cat << EOF > temp.cpp
# // copy the cod here 
# EOF

# COMPILE_ERROR=$(g++ -o temp temp.cpp -DONLINE_JUDGE  -Wl,--stack=268435456 -std=c++20)

# if [ ! -f ./temp ]; then
#     echo "Compilation failed with the following error:"
#     echo "$COMPILE_ERROR"
# else 
#     ulimit -v 254800 

#     OUTPUT=$(echo "100 100" | timeout 1s ./temp 2>&1)
#     EXIT_CODE=$?  # Capture the exit code of the last command

#     if [ $EXIT_CODE -eq 143  ]; then
#         OUTPUT="TLE"
#     elif [ $EXIT_CODE -eq 139 ]; then
#         OUTPUT="Error: Segmentation fault (invalid memory access)."
#     elif [ $EXIT_CODE -eq 136 ]; then
#         OUTPUT="Error: Floating-point exception."
#     elif [ $EXIT_CODE -eq 134 ]; then
#         OUTPUT="Error: Program aborted."
#     elif [ $EXIT_CODE -eq 127 ]; then
#         OUTPUT="Error: Executable not found."
#     elif [ $EXIT_CODE -ge 128 ] && [ $EXIT_CODE -lt 256 ]; then
#         SIGNAL=$((EXIT_CODE - 128))
#         OUTPUT="Error: Program terminated by signal $SIGNAL."
#     elif [ $EXIT_CODE -ne 0 ]; then
#         OUTPUT="Error: Program terminated unexpectedly."
#     fi

#     echo "$OUTPUT , $EXIT_CODE"
# fi

# if [ -f temp ]; then
#     rm temp
# fi

# if [ -f temp.cpp ]; then
#     rm temp.cpp
# fi

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

cat << EOF > temp.cpp
//some code
#include <iostream>
using namespace std ;
int main(){
    int n = 3

    cout<<n*2;
    return ;
}
EOF


OUTPUT=$(g++ -o temp temp.cpp 2>&1)

if [ ! -f ./temp ]; then
    echo "Compilation failed with the following error:"
else 
    ulimit -v 254800 
    #input idhr jaega : echo "$(cat input.txt)"
    OUTPUT=$(echo "2" | timeout 1s ./temp 2>&1)
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
    elif [ $EXIT_CODE -ne 0 ]; then
        OUTPUT="Error: Program terminated unexpectedly."
    fi

fi

if [ -f temp ]; then
    rm temp
fi

if [ -f temp.cpp ]; then
    rm temp.cpp
fi

echo "$OUTPUT "