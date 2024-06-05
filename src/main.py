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
APPLICATION_ID = os.getenv('APPLICATION_ID')

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
    # Check if request is even valid
    if not request.json:
        return jsonify({}), 400
    if request.json['type'] != InteractionType.APPLICATION_COMMAND:
        return jsonify({}), 400
    
    # Prepare response
    data = request.json['data']
    name = data['name']
    
    # Hello command
    if name == 'hello':
        response = {
            'type': InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
            'data': {
                'content': 'Hello, world!'
            }
        }
        return jsonify(response), 200
    
    # Petpet command
    elif name == 'petpet':
        # Fetching options
        options = {option['name']: option['value'] for option in data['options']}

        # Acknowledge the interaction
        interaction_id = request.json['id']
        interaction_token = request.json['token']

        api_url = f'https://discord.com/api/v9/interactions/{interaction_id}/{interaction_token}/callback'
        body = {
            'type': InteractionResponseType.DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE,
            'data': {
                'flags': 64 if options.get('ephemeral') is True else 0
            }
        }
        response = requests.post(api_url, json=body)

        # Get author id
        author_id = None
        if request.json.get('member'):
            author_id = request.json['member']['user']['id']
        elif request.json['user']:
            author_id = request.json['user']['id']

        # Get petpet user
        user_id = options['user']
        avatar_url = None
        resolved_user = data['resolved']['users'][user_id]

        # Get avatar url of petpet user
        if data['resolved'].get('members') and data['resolved']['members'].get(user_id) and (data['resolved']['members'][user_id]['avatar'] is not None):
            avatar_hash = data['resolved']['members'][user_id]['avatar']
            avatar_url = f'https://cdn.discordapp.com/guilds/{request.json['guild_id']}/users/{user_id}/avatars/{avatar_hash}.png?size=1024'
        else:
            avatar_hash = resolved_user['avatar']
            avatar_url = f'https://cdn.discordapp.com/avatars/{user_id}/{avatar_hash}.png?size=1024'
            
        # Generate petpet gif
        if options.get("resolution"):
            petpet.resolution = (options["resolution"], options["resolution"])
        if options.get("frame_delay"):
            petpet.delay = options["frame_delay"]
        if options.get("frame_count"):
            petpet.frames = options["frame_count"]

        avatar_bytes = file_url_to_bytesio(avatar_url)
        output_bytes = BytesIO()

        petpet.make(avatar_bytes, output_bytes)

        # Send the petpet gif
        file_name = 'petpet.gif'
        api_url = f'https://discord.com/api/v9/webhooks/{APPLICATION_ID}/{interaction_token}/messages/@original'
        body = {
            'content': f'<@{author_id}> has pet <@{user_id}>',
            'allowed_mentions': {
                'parse': []
            },
            'attachments': [{ 'id': 0, 'filename': file_name, 'description': f'A gif of a hand patting the avatar of Discord user {resolved_user["global_name"]} ({resolved_user["username"]})' }],
        }
        files = {
            'files[0]': (file_name, output_bytes.getvalue()),
            'payload_json': (None, json.dumps(body), 'application/json')
        }
        api_response = requests.patch(api_url, files=files)
        return jsonify({ 'status': api_response.status_code }), 200

    # Unknown command
    else:
        return jsonify({}), 400
  
@app.route('/')
def index():
  return 'Hello world'

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080, debug=os.getenv('DEBUG') == 'True')