import pytest

from sum_a_b import sum_a_b


@pytest.mark.parametrize(
    "a,b,expected",
    [
        (1, 2, 3),
        (1, 0, 1),
        (0, 0, 0),
        (-1, 0, -1),
        (-1, -1, -2),
    ],
)
def test_sum_a_b_simple(a: float, b: float, expected: float) -> None:
    assert sum_a_b(a, b) == pytest.approx(expected)
