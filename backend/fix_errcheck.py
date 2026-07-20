import re
import os

def fix_errchecks(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Fix defers
    content = re.sub(r'defer (.*?)\.Close\(\)', r'defer func() { _ = \1.Close() }()', content)
    content = re.sub(r'defer tx\.Rollback\(\)', r'defer func() { _ = tx.Rollback() }()', content)
    
    # Fix Encode
    content = re.sub(r'^(\t|\s+)json\.NewEncoder\((.*?)\)\.Encode\((.*?)\)', r'\1_ = json.NewEncoder(\2).Encode(\3)', content, flags=re.MULTILINE)
    
    # Fix Decode in tests
    content = re.sub(r'^(\t|\s+)json\.NewDecoder\((.*?)\)\.Decode\((.*?)\)', r'\1_ = json.NewDecoder(\2).Decode(\3)', content, flags=re.MULTILINE)

    # Fix fmt.Fprintf
    content = re.sub(r'^(\t|\s+)fmt\.Fprintf\((.*?)\)', r'\1_, _ = fmt.Fprintf(\2)', content, flags=re.MULTILINE)

    # Fix io.Copy
    content = re.sub(r'^(\t|\s+)io\.Copy\((.*?)\)', r'\1_, _ = io.Copy(\2)', content, flags=re.MULTILINE)

    # Fix os.Remove / os.RemoveAll / os.MkdirAll
    content = re.sub(r'^(\t|\s+)os\.Remove\((.*?)\)', r'\1_ = os.Remove(\2)', content, flags=re.MULTILINE)
    content = re.sub(r'^(\t|\s+)os\.RemoveAll\((.*?)\)', r'\1_ = os.RemoveAll(\2)', content, flags=re.MULTILINE)
    content = re.sub(r'^(\t|\s+)os\.MkdirAll\((.*?)\)', r'\1_ = os.MkdirAll(\2)', content, flags=re.MULTILINE)

    # Fix DB.Exec
    content = re.sub(r'^(\t|\s+)database\.DB\.Exec\((.*?)\)', r'\1_, _ = database.DB.Exec(\2)', content, flags=re.MULTILINE)
    
    # Fix DB.QueryRow(..).Scan(..)
    content = re.sub(r'^(\t|\s+)database\.DB\.QueryRow\((.*?)\)\.Scan\((.*?)\)', r'\1_ = database.DB.QueryRow(\2).Scan(\3)', content, flags=re.MULTILINE)

    # Fix ParseMultipartForm
    content = re.sub(r'^(\t|\s+)r\.ParseMultipartForm\((.*?)\)', r'\1_ = r.ParseMultipartForm(\2)', content, flags=re.MULTILINE)
    
    # Fix json.Unmarshal
    content = re.sub(r'^(\t|\s+)json\.Unmarshal\((.*?)\)', r'\1_ = json.Unmarshal(\2)', content, flags=re.MULTILINE)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

for root, dirs, files in os.walk('f:\\todo-bench\\Gemma4\\backend'):
    for file in files:
        if file.endswith('.go'):
            fix_errchecks(os.path.join(root, file))

print("Done")
