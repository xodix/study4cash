#!/usr/bin/python3

import json
import sys
import csv
import os
import xml.etree.ElementTree as ET
import pathlib


def get_data_from_csv(file_name: str) -> dict:
    data = []
    with open(file_name, "r") as f:
        reader = csv.reader(f, delimiter=";")
        headers = next(reader)
        print(headers)
        headers = list(map(lambda x: x if len(x.split(";"))
                       == 1 else "Y" + x.split(";")[1], headers))

        for row in reader:
            obj = {}
            for i in range(len(headers)):
                if row[i] == "":
                    continue

                obj[headers[i]] = row[i]

            data.append(obj)

    if os.getenv("DEBUG", False):
        print(data)

    return data


def write_data_to_json(file_name: str):
    with open(file_name, "w") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)


def write_data_to_xml(file_name: str):
    file_prefix = file_name.split(".")[0]
    root = ET.Element("root")

    for item in data:
        record = ET.SubElement(root, file_prefix.lower())
        for key, value in item.items():
            child = ET.SubElement(record, key)
            child.text = str(value)

    tree = ET.ElementTree(root)
    ET.indent(tree)
    tree.write(file_name, encoding="unicode", xml_declaration=True)


if __name__ == "__main__":
    print(sys.argv)
    # get first arg
    if len(sys.argv) < 2:
        raise Exception("Provide the file from witch you want to generate")

    file_name = sys.argv[1]
    if not file_name.endswith(".csv"):
        raise Exception("Input file must end with .csv")

    data = get_data_from_csv(file_name)

    file_prefix = file_name.split(".")[0]
    write_data_to_json(f"{file_prefix}.json")
    write_data_to_xml(f"{file_prefix}.xml")
