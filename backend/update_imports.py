import os
import re

def update_imports(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # We want to replace "todo-backend/pkgName" with "todo-backend/internal/pkgName"
    packages = ["database", "handlers", "middleware", "models", "repository", "service", "testutils"]
    
    for pkg in packages:
        # Match EXACTLY "todo-backend/pkg" to avoid replacing already fixed ones
        content = re.sub(rf'"todo-backend/{pkg}"', f'"todo-backend/internal/{pkg}"', content)
        
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

# Scan cmd and internal directories
for target_dir in ['cmd', 'internal']:
    for root, dirs, files in os.walk(f'f:\\todo-bench\\Gemma4\\backend\\{target_dir}'):
        for file in files:
            if file.endswith('.go'):
                update_imports(os.path.join(root, file))

print("Imports updated!")
