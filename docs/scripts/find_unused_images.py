import os
import re

def main():
    content_dirs = ['docs']  # Add other content directories if needed

    used_images = set()

    for content_dir in content_dirs:
        for root, _, files in os.walk(content_dir):
            md_files = [file for file in files if file.endswith(('.md', '.mdx'))]
            png_files = [file for file in files if file.endswith('.png')]

            if md_files:
                for md_file in md_files:
                    with open(os.path.join(root, md_file), 'r') as content_file:
                        content = content_file.read()
                        used_images.update(re.findall(r'!\[.*?\]\((.*?)\)', content))

                unused_images = [os.path.join(root, png) for png in png_files if png not in used_images]

                if unused_images:
                    print(f'Unused images in {root}:')
                    for img in unused_images:
                        print(img)
                    print()

    if not used_images:
        print('No unused images found.')

if __name__ == '__main__':
    main()