import os
import csv
import re
import argparse
from pathlib import Path

def find_markdown_files(source_dir):
    """Find all markdown files in the source directory and its subdirectories."""
    markdown_files = []
    for root, _, files in os.walk(source_dir):
        for file in files:
            if file.endswith('.md'):
                # Get the relative path from the source directory
                full_path = os.path.join(root, file)
                rel_path = os.path.relpath(full_path, source_dir)
                markdown_files.append(rel_path)
    return markdown_files

def create_csv(markdown_files, csv_path):
    """Create a CSV file with the relative paths and empty URL columns."""
    with open(csv_path, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(['Relative Path', 'Redirect URL'])
        for file in markdown_files:
            writer.writerow([file, ''])
    print(f"CSV created at {csv_path}")

def read_csv(csv_path):
    """Read the CSV file and return a dictionary of file paths to redirect URLs."""
    redirects = {}
    with open(csv_path, 'r', newline='') as csvfile:
        reader = csv.reader(csvfile)
        next(reader)  # Skip header
        for row in reader:
            if len(row) >= 2 and row[1].strip():  # Only include if URL is provided
                redirects[row[0]] = row[1].strip()
    return redirects

def update_markdown_files(source_dir, redirects):
    """Update markdown files with redirect metadata."""
    for rel_path, redirect_url in redirects.items():
        full_path = os.path.join(source_dir, rel_path)
        
        if not os.path.exists(full_path):
            print(f"Warning: File {full_path} does not exist. Skipping.")
            continue
            
        with open(full_path, 'r', encoding='utf-8') as file:
            content = file.read()
        
        # Check if the file already has a metadata block
        metadata_pattern = r'^---\s*\n(.*?)\n---\s*\n'
        metadata_match = re.search(metadata_pattern, content, re.DOTALL)
        
        if metadata_match:
            # File has a metadata block, check if it already has a redirect
            metadata = metadata_match.group(1)
            if 'redirect:' in metadata:
                # Replace existing redirect
                updated_metadata = re.sub(
                    r'redirect:.*', 
                    f'redirect: {redirect_url}', 
                    metadata
                )
            else:
                # Add redirect to existing metadata
                updated_metadata = metadata + f"\nredirect: {redirect_url}"
            
            # Replace the old metadata with updated metadata
            updated_content = re.sub(
                metadata_pattern,
                f'---\n{updated_metadata}\n---\n',
                content,
                count=1
            )
        else:
            # No metadata block, create a new one
            updated_content = f'---\nredirect: {redirect_url}\n---\n{content}'
        
        # Write the updated content back to the file
        with open(full_path, 'w', encoding='utf-8') as file:
            file.write(updated_content)
        
        print(f"Updated {rel_path} with redirect to {redirect_url}")

def main():
    parser = argparse.ArgumentParser(description='MkDocs to AWS Docs Migration Tool')
    parser.add_argument('--source', required=True, help='Source directory containing markdown files')
    parser.add_argument('--csv', required=True, help='Path to the CSV file (to create or to read)')
    parser.add_argument('--mode', choices=['create_csv', 'update_files'], required=True, 
                        help='Mode: create CSV or update files with redirects')
    
    args = parser.parse_args()
    
    # Create absolute paths
    source_dir = os.path.abspath(args.source)
    csv_path = os.path.abspath(args.csv)
    
    if args.mode == 'create_csv':
        markdown_files = find_markdown_files(source_dir)
        create_csv(markdown_files, csv_path)
        print(f"Found {len(markdown_files)} markdown files. Edit the CSV at {csv_path} to add redirect URLs.")
    
    elif args.mode == 'update_files':
        if not os.path.exists(csv_path):
            print(f"Error: CSV file {csv_path} does not exist.")
            return
        
        redirects = read_csv(csv_path)
        if not redirects:
            print("No redirects found in the CSV or all redirect URLs are empty.")
            return
            
        update_markdown_files(source_dir, redirects)
        print(f"Updated {len(redirects)} markdown files with redirects.")

if __name__ == "__main__":
    main()