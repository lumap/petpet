from typing import Literal
import requests
from io import BytesIO

def get_image_from_url(file_url: str) -> BytesIO | Literal[False]:
    try:
        response = requests.get(file_url)
        response.raise_for_status()
        file_bytes = BytesIO(response.content)
        return file_bytes
    except requests.RequestException:
        return False
