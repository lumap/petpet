from mimetypes import guess_type, add_type
import os
import dotenv

add_type('image/webp', '.webp')

from src.utils.interactions import defer_interaction, finish_interaction, finish_interaction_upload_img, reply_early_to_interaction
from src.utils.make_petpet_gif import make_petpet_gif
from src.utils.get_image_from_url import get_image_from_url

if os.path.isfile(".env.dev"):
    dotenv.load_dotenv(dotenv.find_dotenv(filename='.env.dev'))
else:
    dotenv.load_dotenv(dotenv.find_dotenv(filename='.env.prod'))

from flask import Flask, request, jsonify, g
from discord_interactions import verify_key_decorator, InteractionType, InteractionResponseType
import uuid
import time


CLIENT_PUBLIC_KEY = os.getenv('CLIENT_PUBLIC_KEY')
APPLICATION_ID = os.getenv('APPLICATION_ID')

blacklisted_users = ["1118171067463241868"]

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
    # Python type checking can be stoopid
    if not APPLICATION_ID:
        return jsonify({}), 500
    
    # Check if request is even valid
    if not request.json:
        return jsonify({}), 400
    if request.json['type'] != InteractionType.APPLICATION_COMMAND:
        return jsonify({}), 400
    
    # Prepare response
    data = request.json['data']
    name = data['name']
    
    # Hello command
    if name == 'meow':
        response = {
            'type': InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
            'data': {
                'content': 'Meow, meow meow~!'
            }
        }
        return jsonify(response), 200
    
    # Petpet command
    elif name == 'petpet':
        
        # Divide code into subcommands
        subcommand = data['options'][0]
        options = {option['name']: option['value'] for option in subcommand['options']}
        
        # user 
        if subcommand['name'] == 'user':
            
            # Acknowledge the interaction
            interaction_id = request.json['id']
            interaction_token = request.json['token']


            # Get petpet user
            user_id = options['user']
            
            if user_id in blacklisted_users:
                return jsonify({
                    'type': 4,
                    'data': {
                        'content': 'This user is blacklisted.',
                        'flags': 64
                    }
                })

            avatar_url = None
            resolved_user = data['resolved']['users'][user_id]

            # Get avatar url of petpet user
            if data['resolved'].get('members') and data['resolved']['members'].get(user_id) and (data['resolved']['members'][user_id]['avatar'] is not None) and (options.get("use_server_avatar", True)):
                avatar_hash = data['resolved']['members'][user_id]['avatar']
                avatar_url = f'https://cdn.discordapp.com/guilds/{request.json['guild_id']}/users/{user_id}/avatars/{avatar_hash}.png?size=1024'
            else:
                avatar_hash = resolved_user['avatar']
                if avatar_hash:
                    avatar_url = f'https://cdn.discordapp.com/avatars/{user_id}/{avatar_hash}.png?size=1024'
                else:
                    discrim = resolved_user.get('discriminator')
                    print(discrim)
                    index = ((int(user_id) >> 22) % 6) if discrim == "0" else (int(discrim) % 5)
                    avatar_url = f'https://cdn.discordapp.com/embed/avatars/{index}.png'
                    
            # Get image
            defer_interaction(id=interaction_id, token=interaction_token, ephemeral=options.get('ephemeral', False))
            
            
            image_bytes = get_image_from_url(avatar_url)
            if image_bytes is False:
                finish_interaction(token=interaction_token, msg_content="There was an error fetching the image. Please try again.", app_id=APPLICATION_ID)
                return jsonify({}), 400
            
            petpet = make_petpet_gif(bytes=image_bytes, width=options.get("width", 128), height=options.get("height", 128), speed=options.get("speed", 1))
            
            # Send the petpet gif
            author_id = None
            if request.json.get('member'):
                author_id = request.json['member']['user']['id']
            elif request.json['user']:
                author_id = request.json['user']['id']
            
            attachment_alt_text = f'A gif of a hand patting the avatar of Discord user { f"{resolved_user["global_name"]} ({resolved_user["username"]})" if resolved_user.get("global_name") else resolved_user["username"] }'
            msg_content = f'<@{author_id}> has pet <@{user_id}>'
            
            status_code = finish_interaction_upload_img(petpet=petpet, token=interaction_token, msg_content=msg_content, attachment_alt_text=attachment_alt_text, app_id=APPLICATION_ID)

            return jsonify({ 'status': status_code }), 200
        
        # image_via_url
        elif subcommand['name'] == 'image_via_url':
            
            # Acknowledge the interaction
            interaction_id = request.json['id']
            interaction_token = request.json['token']

            url = options['image_url']
            guessed_type = guess_type(url)[0]
            if not guessed_type or guessed_type.split('/')[0] != 'image':
                reply_early_to_interaction(id=interaction_id, token=interaction_token, content="The provided URL is invalid.")
                return jsonify({}), 400
            
            # Get image
            defer_interaction(id=interaction_id, token=interaction_token, ephemeral=options.get('ephemeral', False))
            
            image_bytes = get_image_from_url(url)
            if image_bytes is False:
                finish_interaction(token=interaction_token, msg_content="There was an error fetching the image. Please try again.", app_id=APPLICATION_ID)
                return jsonify({}), 400
            
            petpet = make_petpet_gif(bytes=image_bytes, width=options.get("width", 128), height=options.get("height", 128), speed=options.get("speed", 1))
            
            # Send the petpet gif
            author_id = None
            if request.json.get('member'):
                author_id = request.json['member']['user']['id']
            elif request.json['user']:
                author_id = request.json['user']['id']
            
            attachment_alt_text = f'A gif of a hand patting the avatar of a user via URL'
            msg_content = f'<@{author_id}> has pet an image via URL'
            
            status_code = finish_interaction_upload_img(petpet=petpet, token=interaction_token, msg_content=msg_content, attachment_alt_text=attachment_alt_text, app_id=APPLICATION_ID)

            return jsonify({ 'status': status_code }), 200
        
        # image_via_upload
        elif subcommand['name'] == 'image_via_upload':
            
            # Acknowledge the interaction
            interaction_id = request.json['id']
            interaction_token = request.json['token']
            
            image_id = options['image_upload']
            image = data['resolved']['attachments'][image_id]
            
            if not image["content_type"].startswith("image/"):
                reply_early_to_interaction(id=interaction_id, token=interaction_token, content="The uploaded file is not an image.")
                return jsonify({}), 400
            
            # Get image
            defer_interaction(id=interaction_id, token=interaction_token, ephemeral=options.get('ephemeral', False))
            
            image_bytes = get_image_from_url(image['url'])
            if image_bytes is False:
                finish_interaction(token=interaction_token, msg_content="There was an error fetching the image. Please try again.", app_id=APPLICATION_ID)
                return jsonify({}), 400
            
            petpet = make_petpet_gif(bytes=image_bytes, width=options.get("width", 128), height=options.get("height", 128), speed=options.get("speed", 1))
            
            # Send the petpet gif
            author_id = None
            if request.json.get('member'):
                author_id = request.json['member']['user']['id']
            elif request.json['user']:
                author_id = request.json['user']['id']
            
            attachment_alt_text = f'A gif of a hand patting the avatar of a user via upload'
            msg_content = f'<@{author_id}> has pet an image via upload'
            
            status_code = finish_interaction_upload_img(petpet=petpet, token=interaction_token, msg_content=msg_content, attachment_alt_text=attachment_alt_text, app_id=APPLICATION_ID)

            return jsonify({ 'status': status_code }), 200

        # Unknown subcommand
        else:
            return jsonify({}), 400

    # petpet user command
    elif name == 'PetPet this user':
        # Acknowledge the interaction
        interaction_id = request.json['id']
        interaction_token = request.json['token']

        # Get petpet user
        user_id = data['target_id']
        
        if user_id in blacklisted_users:
            return jsonify({
                'type': 4,
                'data': {
                    'content': 'This user is blacklisted.',
                    'flags': 64
                }
            })

        avatar_url = None
        resolved_user = data['resolved']['users'][user_id]

        # Get avatar url of petpet user
        if data['resolved'].get('members') and data['resolved']['members'].get(user_id) and (data['resolved']['members'][user_id]['avatar'] is not None):
            avatar_hash = data['resolved']['members'][user_id]['avatar']
            avatar_url = f'https://cdn.discordapp.com/guilds/{request.json['guild_id']}/users/{user_id}/avatars/{avatar_hash}.png?size=1024'
        else:
            avatar_hash = resolved_user['avatar']
            if avatar_hash:
                avatar_url = f'https://cdn.discordapp.com/avatars/{user_id}/{avatar_hash}.png?size=1024'
            else:
                discrim = resolved_user.get('discriminator')
                index = ((int(user_id) >> 22) % 6) if discrim is None else (int(discrim) % 5)
                avatar_url = f'https://cdn.discordapp.com/embed/avatars/{index}.png'
                
        # Get image
        defer_interaction(id=interaction_id, token=interaction_token)

        image_bytes = get_image_from_url(avatar_url)
        if image_bytes is False:
            finish_interaction(token=interaction_token, msg_content="There was an error fetching the image. Please try again.", app_id=APPLICATION_ID)
            return jsonify({}), 400
    
        petpet = make_petpet_gif(bytes=image_bytes)
        
        # Send the petpet gif
        author_id = None
        if request.json.get('member'):
            author_id = request.json['member']['user']['id']
        elif request.json['user']:
            author_id = request.json['user']['id']
        
        attachment_alt_text = f'A gif of a hand patting the avatar of Discord user {resolved_user["global_name"]} ({resolved_user["username"]})'
        msg_content = f'<@{author_id}> has pet <@{user_id}>'
        
        status_code = finish_interaction_upload_img(petpet=petpet, token=interaction_token, msg_content=msg_content, attachment_alt_text=attachment_alt_text, app_id=APPLICATION_ID)

        return jsonify({ 'status': status_code }), 200
    
    # petpet message command
    elif name == 'PetPet this message\'s author':
        
        # Acknowledge the interaction
        interaction_id = request.json['id']
        interaction_token = request.json['token']
                
        # Get message author
        message_target_id = data['target_id']
        resolved_user = data['resolved']['messages'][message_target_id]['author']
        
        # get petpet user
        user_id = resolved_user['id']
        
        if user_id in blacklisted_users:
            return jsonify({
                'type': 4,
                'data': {
                    'content': 'This user is blacklisted.',
                    'flags': 64
                }
            })        
        avatar_url = None

        # Get avatar url of petpet user
        if data['resolved'].get('members') and data['resolved']['members'].get(user_id) and (data['resolved']['members'][user_id]['avatar'] is not None):
            avatar_hash = data['resolved']['members'][user_id]['avatar']
            avatar_url = f'https://cdn.discordapp.com/guilds/{request.json['guild_id']}/users/{user_id}/avatars/{avatar_hash}.png?size=1024'
        else:
            avatar_hash = resolved_user['avatar']
            if avatar_hash:
                avatar_url = f'https://cdn.discordapp.com/avatars/{user_id}/{avatar_hash}.png?size=1024'
            else:
                discrim = resolved_user.get('discriminator')
                index = ((int(user_id) >> 22) % 6) if discrim is None else (int(discrim) % 5)
                avatar_url = f'https://cdn.discordapp.com/embed/avatars/{index}.png'
                
        # Get image
        defer_interaction(id=interaction_id, token=interaction_token)
        
        image_bytes = get_image_from_url(avatar_url)
        if image_bytes is False:
            finish_interaction(token=interaction_token, msg_content="There was an error fetching the image. Please try again.", app_id=APPLICATION_ID)
            return jsonify({}), 400

        petpet = make_petpet_gif(bytes=image_bytes)
        
        # Send the petpet gif
        author_id = None
        if request.json.get('member'):
            author_id = request.json['member']['user']['id']
        elif request.json['user']:
            author_id = request.json['user']['id']
        
        attachment_alt_text = f'A gif of a hand patting the avatar of Discord user {resolved_user["global_name"]} ({resolved_user["username"]})'
        msg_content = f'<@{author_id}> has pet <@{user_id}>'
        
        status_code = finish_interaction_upload_img(petpet=petpet, token=interaction_token, msg_content=msg_content, attachment_alt_text=attachment_alt_text, app_id=APPLICATION_ID)

        return jsonify({ 'status': status_code }), 200
        return jsonify({}), 200
    
    # Unknown command type
    else:
        return jsonify({}), 400
  
@app.route('/')
def index():
  return 'Hello world'

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8080, debug=os.getenv('DEBUG') == 'True')