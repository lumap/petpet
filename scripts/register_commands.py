#!.venv/bin/python3

import dotenv
dotenv.load_dotenv(dotenv.find_dotenv(filename='.env.dev'))

import requests
import os
import argparse

def register_commands(json_path: str, bot_token: str, client_id: str, guild_id: str = None):
    headers = {
        'Authorization': f'Bot {bot_token}',
        'Content-Type': 'application/json'
    }
    url = f'https://discord.com/api/v9/applications/{client_id}/commands'
    if guild_id:
        url = f'https://discord.com/api/v9/applications/{client_id}/guilds/{guild_id}/commands'
    with open(json_path, 'r') as file:
        commands = file.read()
        response = requests.put(url, headers=headers, data=commands)
        try:
            print(response.json())
        except:
            print(response.text)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Register commands to Discord API')
    parser.add_argument('json_path', type=str, help='Path to the JSON file containing the commands')

    args = parser.parse_args()
    json_path = args.json_path
    register_commands(json_path, os.getenv('BOT_TOKEN'), os.getenv('CLIENT_ID'), os.getenv('GUILD_ID'))