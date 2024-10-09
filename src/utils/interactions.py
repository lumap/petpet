import requests
from discord_interactions import InteractionResponseType
import json

def reply_early_to_interaction(id: str, token: str, content: str) -> None:
    api_url = f'https://discord.com/api/v9/interactions/{id}/{token}/callback'
    body = {
        'type': InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
        'data': {
            'content': content,
            'flags': 64
        }
    }
    requests.post(api_url, json=body)

def defer_interaction(id: str, token: str, ephemeral: bool = False) -> None:
    api_url = f'https://discord.com/api/v9/interactions/{id}/{token}/callback'
    body = {
        'type': InteractionResponseType.DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE,
        'data': {
            'flags': 64 if ephemeral else 0
        }
    }
    requests.post(api_url, json=body)
    
def finish_interaction(token: str, msg_content: str, app_id: str) -> int:
    api_url = f'https://discord.com/api/v9/webhooks/{app_id}/{token}/messages/@original'
    body = {
        'content': msg_content
    }
    api_response = requests.patch(api_url, json=body)
    return api_response.status_code
    
    
def finish_interaction_upload_img(petpet: bytes, token: str, msg_content: str, attachment_alt_text: str, app_id: str) -> int:
    file_name = 'petpet.gif'
    api_url = f'https://discord.com/api/v9/webhooks/{app_id}/{token}/messages/@original'
    body = {
        'content': msg_content,
        'allowed_mentions': {
            'parse': []
        },
        'attachments': [{ 'id': 0, 'filename': file_name, 'description': attachment_alt_text }],
    }
    files = {
        'files[0]': (file_name, petpet),
        'payload_json': (None, json.dumps(body), 'application/json')
    }
    api_response = requests.patch(api_url, files=files)
    return api_response.status_code