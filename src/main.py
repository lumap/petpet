import dotenv
dotenv.load_dotenv()
dotenv.load_dotenv(dotenv.find_dotenv(filename='.env.dev'))

import os
from flask import Flask, request, jsonify, g
from discord_interactions import verify_key_decorator, InteractionType, InteractionResponseType
from petpetgif import petpet
import requests
from io import BytesIO
import json
import uuid
import time

CLIENT_PUBLIC_KEY = os.getenv('CLIENT_PUBLIC_KEY')

def file_url_to_bytesio(file_url):
    response = requests.get(file_url)
    response.raise_for_status()  # Check if the request was successful
    file_bytes = BytesIO(response.content)
    return file_bytes

app = Flask(__name__)

@app.before_request
def before_request_func():
    execution_id = uuid.uuid4()
    g.start_time = time.time()
    g.execution_id = execution_id

    print(g.execution_id, "ROUTE CALLED ", request.url)

@app.route('/interactions', methods=['POST'])
@verify_key_decorator(CLIENT_PUBLIC_KEY)
def interactions():
  if request.json['type'] == InteractionType.APPLICATION_COMMAND:
    data = request.json['data']
    name = data['name']
    response = {
        'type': InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
        'data': {
            'content': 'Uhhh'
        }
    }
    if name == 'hello':
        response['data']['content'] = 'Hello world'
        return jsonify(response)
    elif name == 'petpet':
        author_id = None
        if request.json.get('member'):
            author_id = request.json['member']['user']['id']
        elif request.json['user']:
            author_id = request.json['user']['id']

        options = {option['name']: option['value'] for option in data['options']}

        user_id = options['user']
        avatar_url = None
        resolved_user = data['resolved']['users'][user_id]
        if data['resolved'].get('members') and data['resolved']['members'].get(user_id) and (data['resolved']['members'][user_id]['avatar'] is not None):
            avatar_hash = data['resolved']['members'][user_id]['avatar']
            avatar_url = f'https://cdn.discordapp.com/guilds/{request.json['guild_id']}/users/{user_id}/avatars/{avatar_hash}.png?size=1024'
        else:
            avatar_hash = resolved_user['avatar']
            avatar_url = f'https://cdn.discordapp.com/avatars/{user_id}/{avatar_hash}.png?size=1024'
        avatar_bytes = file_url_to_bytesio(avatar_url)
        output_bytes = BytesIO()

        if options.get("resolution"):
            petpet.resolution = (options["resolution"], options["resolution"])
        if options.get("frame_delay"):
            petpet.delay = options["frame_delay"]
        if options.get("frame_count"):
            petpet.frames = options["frame_count"]
        petpet.make(avatar_bytes, output_bytes)
        
        file_name = 'petpet.gif'
        response['data'] = {
            'content': f'<@{author_id}> has pet <@{user_id}>',
            'allowed_mentions': {
                'parse': []
            },
            'attachments': [{ 'id': 0, 'filename': file_name, 'description': f'A gif of a hand patting the avatar of Discord user {resolved_user["global_name"]} ({resolved_user["username"]})' }],
            'flags': 64 if options.get('ephemeral') is True else 0
        }
        files = {
            'files[0]': (file_name, output_bytes.getvalue()),
            'payload_json': (None, json.dumps(response), 'application/json')
        }
        print(response)
        api_url = f'https://discord.com/api/v9/interactions/{request.json["id"]}/{request.json["token"]}/callback'
        api_response = requests.post(api_url, files=files)
        try:
           print(api_response.json())
        except:
            print(api_response.text)

        return api_response.text, api_response.status_code
    else:
        return jsonify(response)

  
@app.route('/')
def index():
  return 'Hello world'

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080, debug=os.getenv('DEBUG', False))