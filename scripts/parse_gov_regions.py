import json
import sys

def transform_and_clean(input_file):
    with open(input_file, 'r', encoding='utf-8') as f:
        data = json.load(f)

    seen_names = set() 
    transformed_data = []
    regions = []         
    departments = []     

    for entry in data:
        region_name = entry.get('region_name')
        if region_name and region_name not in seen_names:
            seen_names.add(region_name)
            regions.append({'Name': region_name, 'Country': 'France'})

        dep_name = entry.get('dep_name')
        if dep_name and dep_name not in seen_names:
            seen_names.add(dep_name)
            departments.append({'Name': dep_name, 'Country': 'France'})

    transformed_data = {
        'regions': regions,
        'departments': departments
    }

    output_file = input_file + '.prepared'

    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("// Regions\n\n")
        for entry in transformed_data['regions']:
            f.write(f"{{Name: \"{entry['Name']}\", Country: \"{entry['Country']}\"}},\n")
        
        f.write("\n// Departments\n\n")
        for entry in transformed_data['departments']:
            f.write(f"{{Name: \"{entry['Name']}\", Country: \"{entry['Country']}\"}},\n")

    print(f"Transformation complete. Output saved to: {output_file}")
    print(f"Total regions extracted: {len(regions)}")
    print(f"Total departments extracted: {len(departments)}")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python parse_gov_regions.py <input_file>")
        sys.exit(1)

    input_file = sys.argv[1]

    transform_and_clean(input_file)
