import os
import re
import zipfile

def clean_go_code(text):
    def replacer(match):
        s = match.group(0)

        if s.startswith('/'):
            return ""

        return s

    pattern = re.compile(r'//.*?$|/\*.*?\*/|"(?:\\.|[^\\"])*"|`(?:[^`])*`', re.DOTALL | re.MULTILINE)
    return re.sub(pattern, replacer, text)

def clean_python_code(text):
    def replacer(match):
        s = match.group(0)

        if s.startswith('#'):
            return ""

        return s

    pattern = re.compile(r'#.*?$|"""(?:\\.|[^"\\])*"""|\'\'\'(?:\\.|[^\'\\])*\'\'\'|"(?:\\.|[^\\"])*"|\'(?:\\.|[^\'\\])*\'', re.DOTALL | re.MULTILINE)

    cleaned = re.sub(pattern, replacer, text)

    cleaned = re.sub(r'^\s*$', '', cleaned, flags=re.MULTILINE)
    return cleaned

def main():
    zip_filename = "clean_project.zip"

    with zipfile.ZipFile(zip_filename, 'w', zipfile.ZIP_DEFLATED) as zipf:
        for root, dirs, files in os.walk("."):

            if "vendor" in dirs:
                dirs.remove("vendor")
            if ".git" in dirs:
                dirs.remove(".git")
            if "__pycache__" in dirs:
                dirs.remove("__pycache__")

            for file in files:

                if file == zip_filename or file == "cleaner.py":
                    continue

                file_path = os.path.join(root, file)
                arcname = os.path.relpath(file_path, ".")

                try:
                    if file.endswith(".go"):
                        with open(file_path, "r", encoding="utf-8") as f:
                            content = f.read()
                        cleaned = clean_go_code(content)
                        zipf.writestr(arcname, cleaned)

                    elif file.endswith(".py"):
                        with open(file_path, "r", encoding="utf-8") as f:
                            content = f.read()
                        cleaned = clean_python_code(content)
                        zipf.writestr(arcname, cleaned)

                    else:

                        zipf.write(file_path, arcname)

                except Exception as e:
                    print(f"Error processing {file_path}: {e}")
                    zipf.write(file_path, arcname)

    print(f"✅ Project successfully cleaned and saved to: {zip_filename}")

if __name__ == "__main__":
    main()