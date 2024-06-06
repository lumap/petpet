import requests
from io import BytesIO

def file_url_to_bytesio(file_url: str) -> BytesIO:
    response = requests.get(file_url)
    response.raise_for_status()
    file_bytes = BytesIO(response.content)
    return file_bytes