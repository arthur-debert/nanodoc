import os

from nanodoc.core import process_file
from nanodoc.files import create_content_item


def test_process_file_no_line_numbers(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = process_file(create_content_item(file_path), None, 0)
    assert "Line 1" in output
    assert "Line 2" in output
    assert "1:" not in output
    assert "2:" not in output


def test_process_file_with_line_numbers_all(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = process_file(create_content_item(file_path), "all", 0)
    assert "1: Line 1" in output
    assert "2: Line 2" in output


def test_process_file_with_line_numbers_file(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = process_file(create_content_item(file_path), "file", 0)
    assert "1: Line 1" in output
    assert "2: Line 2" in output


def test_process_file_not_found():
    file_path = "nonexistent_file.txt"
    output, _ = process_file(create_content_item(file_path), None, 0)
    assert "Error: File not found" in output


def test_process_file_header_assignment(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    file_path = str(test_file)
    output, _ = process_file(create_content_item(file_path), None, 0)
    header = "\n" + os.path.basename(file_path) + "\n\n"
    assert output.startswith(header)


def test_process_file_with_no_header(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")
    file_path = str(test_file)
    output, _ = process_file(create_content_item(file_path), None, 0, show_header=False)
    assert "Line 1" in output
    assert "Line 2" in output
    assert not output.startswith("\n")
    assert os.path.basename(file_path) not in output.split("\n")[0]


def test_process_file_with_header_sequence(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    file_path = str(test_file)

    # Test numerical sequence
    output, _ = process_file(
        create_content_item(file_path), None, 0, sequence="numerical", seq_index=0
    )
    assert "\n1. test_file.txt\n\n" in output

    # Test letter sequence
    output, _ = process_file(
        create_content_item(file_path), None, 0, sequence="letter", seq_index=1
    )
    assert "\nb. test_file.txt\n\n" in output

    # Test roman sequence
    output, _ = process_file(
        create_content_item(file_path), None, 0, sequence="roman", seq_index=2
    )
    assert "\niii. test_file.txt\n\n" in output


def test_process_file_with_header_style(tmpdir):
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    file_path = str(test_file)

    # Test nice style with sequence
    output, _ = process_file(
        create_content_item(file_path),
        None,
        0,
        sequence="numerical",
        seq_index=0,
        style="nice",
    )
    assert "\n1. Test File (test_file.txt)\n\n" in output
