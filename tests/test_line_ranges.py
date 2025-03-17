import pytest

from nanodoc.data import LineRange, validate_content_item
from nanodoc.files import create_content_item, get_file_content, parse_line_reference


def test_parse_line_reference_single_line():
    # Test parsing a single line reference
    result = parse_line_reference("L10")
    assert len(result) == 1
    assert result[0].start == 10
    assert result[0].end == 10


def test_parse_line_reference_range():
    # Test parsing a range reference
    result = parse_line_reference("L5-10")
    assert len(result) == 1
    assert result[0].start == 5
    assert result[0].end == 10


def test_parse_line_reference_multiple():
    # Test parsing multiple references
    result = parse_line_reference("L5,L10-15,L20")
    assert len(result) == 3
    assert [r.start for r in result] == [5, 10, 20]
    assert [r.end for r in result] == [5, 15, 20]


def test_parse_line_reference_invalid():
    # Test parsing invalid references
    with pytest.raises(ValueError):
        parse_line_reference("L")  # Empty line number

    with pytest.raises(ValueError):
        parse_line_reference("L-5")  # Missing start number

    with pytest.raises(ValueError):
        parse_line_reference("L5-")  # Missing end number

    with pytest.raises(ValueError):
        parse_line_reference("L5-3")  # End less than start

    with pytest.raises(ValueError):
        parse_line_reference("L5bad")  # Invalid characters after number

    with pytest.raises(ValueError):
        parse_line_reference("L5-10bad")  # Invalid characters after range

    with pytest.raises(ValueError):
        parse_line_reference("L5,L10bad")  # Invalid characters in second part


def test_get_file_content_entire_file(tmpdir):
    # Test getting the entire file content
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    content = get_file_content(file_path)
    assert content == "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"


def test_get_file_content_single_line(tmpdir):
    # Test getting a single line
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    content = get_file_content(file_path, line=3)
    assert content == "Line 3"


def test_get_file_content_line_range(tmpdir):
    # Test getting a range of lines
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    content = get_file_content(file_path, start=2, end=4)
    assert content == "Line 2\nLine 3\nLine 4"


def test_get_file_content_multiple_parts(tmpdir):
    # Test getting multiple parts
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    parts = [LineRange(1, 1), LineRange(3, 4)]
    content = get_file_content(file_path, parts=parts)
    assert content == "Line 1\nLine 3\nLine 4"


def test_get_file_content_line_out_of_range(tmpdir):
    # Test getting a line that is out of range
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3")
    file_path = str(test_file)

    with pytest.raises(ValueError) as excinfo:
        get_file_content(file_path, line=10)
    assert "Line 10 out of range (1-3)" in str(excinfo.value)


def test_verify_path_with_line_reference_valid(tmpdir):
    # Test verifying a path with a valid line reference
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
    file_path = str(test_file)

    # Valid single line
    content_item = create_content_item(f"{file_path}:L3")
    result = content_item.file_path
    line_parts = content_item.ranges
    assert result == file_path
    assert len(line_parts) == 1
    assert line_parts[0].start == 3
    assert line_parts[0].end == 3

    # Valid range
    content_item = create_content_item(f"{file_path}:L2-4")
    result = content_item.file_path
    line_parts = content_item.ranges
    assert result == file_path
    assert len(line_parts) == 1
    assert line_parts[0].start == 2
    assert line_parts[0].end == 4

    # Valid multiple
    content_item = create_content_item(f"{file_path}:L1,L3-4,L5")
    assert len(content_item.ranges) == 3


def test_verify_path_with_line_reference_invalid(tmpdir):
    # Test verifying a path with an invalid line reference
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2\nLine 3")
    file_path = str(test_file)

    # Invalid line number (too high)
    with pytest.raises(ValueError) as excinfo:
        validate_content_item(create_content_item(f"{file_path}:L10"))
    assert "Line reference out of range" in str(excinfo.value)

    # Invalid range (end too high)
    with pytest.raises(ValueError) as excinfo:
        validate_content_item(create_content_item(f"{file_path}:L2-10"))
    assert "Line reference out of range" in str(excinfo.value)

    # Invalid multiple (one invalid)
    with pytest.raises(ValueError) as excinfo:
        validate_content_item(create_content_item(f"{file_path}:L1,L10"))
    assert "Line reference out of range" in str(excinfo.value)

    # Invalid characters after line number
    with pytest.raises(ValueError) as excinfo:
        create_content_item(f"{file_path}:L3bad")
    assert "Invalid character in line reference" in str(excinfo.value)

    # Invalid characters after range
    with pytest.raises(ValueError) as excinfo:
        create_content_item(f"{file_path}:L1-2bad")
    assert "Invalid character in line reference" in str(excinfo.value)

    # Non-existent file with line reference
    with pytest.raises(FileNotFoundError) as excinfo:
        validate_content_item(create_content_item("nonexistent_file.txt:L5"))
    assert "File not found" in str(excinfo.value)
