import re
import sys

def find_main_class(java_code):
    class_pattern = re.compile(r'\bclass\s+(\w+)\b')
    main_pattern = re.compile(r'\bpublic\s+static\s+void\s+main\s*\(\s*String\s*\[\s*\]\s*\w*\s*\)')

    classes = class_pattern.findall(java_code)
    for class_name in classes:
        class_start = java_code.find(f'class {class_name}')
        brace_count = 0
        class_end = class_start
        for i in range(class_start, len(java_code)):
            if java_code[i] == '{':
                brace_count += 1
            elif (java_code[i] == '}') and (brace_count > 0):
                brace_count -= 1
                if brace_count == 0:
                    class_end = i
                    break
        class_body = java_code[class_start:class_end+1]
        if main_pattern.search(class_body):
            return class_name
    return None

if __name__ == "__main__":
    java_code = sys.argv[1]
    class_name = find_main_class(java_code)
    if class_name:
        print(class_name)
    else:
        print("No main method found")