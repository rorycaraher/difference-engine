import os
from unittest.mock import patch, MagicMock
import pytest

from mixer.mixer import Mixer


@pytest.fixture
def mixer():
    return Mixer()


def test_get_stems_returns_paths(mixer, tmp_path):
    paths = mixer.get_stems(str(tmp_path), [1, 2, 3])
    assert len(paths) == 3
    assert all(p.endswith(".mp3") for p in paths)
    assert all(str(tmp_path) in p for p in paths)


def test_get_stems_includes_stem_name(mixer, tmp_path):
    paths = mixer.get_stems(str(tmp_path), ["kick", "bass"])
    assert any("kick" in p for p in paths)
    assert any("bass" in p for p in paths)


@patch("mixer.mixer.ffmpeg.run")
@patch("mixer.mixer.ffmpeg.filter")
@patch("mixer.mixer.ffmpeg.input")
def test_create_mixdown_calls_ffmpeg(mock_input, mock_filter, mock_run, mixer, tmp_path, monkeypatch):
    monkeypatch.setattr("mixer.mixer.OUTPUT_DIR", str(tmp_path))

    # chain: input().filter() returns a mock; ffmpeg.filter().output() returns a mock
    mock_stream = MagicMock()
    mock_input.return_value.filter.return_value = mock_stream
    mock_combined = MagicMock()
    mock_filter.return_value.output.return_value = mock_combined

    output = mixer.create_mixdown(["a.mp3", "b.mp3"], [0.8, 0.6])

    mock_run.assert_called_once_with(mock_combined)
    assert output.endswith(".mp3")
    assert "output-" in output


@patch("mixer.mixer.ffmpeg.run")
@patch("mixer.mixer.ffmpeg.filter")
@patch("mixer.mixer.ffmpeg.input")
def test_create_mixdown_default_volumes(mock_input, mock_filter, mock_run, mixer, tmp_path, monkeypatch):
    monkeypatch.setattr("mixer.mixer.OUTPUT_DIR", str(tmp_path))
    mock_stream = MagicMock()
    mock_input.return_value.filter.return_value = mock_stream
    mock_filter.return_value.output.return_value = MagicMock()

    # called twice with same Mixer instance to confirm default volumes aren't shared
    mixer.create_mixdown(["a.mp3", "b.mp3"])
    mixer.create_mixdown(["a.mp3", "b.mp3"])
    assert mock_run.call_count == 2
