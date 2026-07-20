import os
import re

def add_doc_comments(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        lines = f.readlines()
        
    out = []
    i = 0
    while i < len(lines):
        line = lines[i]
        
        # Match exported types, funcs, vars, constants
        match = re.match(r'^(type|func|var|const)\s+([A-Z]\w*)', line)
        
        # Special match for funcs with receivers: func (r *Receiver) ExportedFunc()
        func_match = re.match(r'^func\s+\([^)]+\)\s+([A-Z]\w*)', line)
        
        symbol_name = None
        if match:
            symbol_name = match.group(2)
        elif func_match:
            symbol_name = func_match.group(1)
            
        if symbol_name:
            # Check if previous line is a comment
            has_comment = False
            if i > 0 and lines[i-1].strip().startswith('//'):
                has_comment = True
            
            if not has_comment:
                # Add doc comment
                out.append(f"// {symbol_name} ...\n")
                
        out.append(line)
        i += 1
        
    with open(filepath, 'w', encoding='utf-8') as f:
        f.writelines(out)

for root, dirs, files in os.walk('f:\\todo-bench\\Gemma4\\backend'):
    for file in files:
        if file.endswith('.go'):
            add_doc_comments(os.path.join(root, file))

print("Doc comments added!")
