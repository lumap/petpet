#!.venv/bin/python3

import requests
import argparse

def register_commands(json_path: str, bot_token: str, client_id: str):
    headers = {
        'Authorization': f'Bot {bot_token}',
        'Content-Type': 'application/json'
    }
    url = f'https://discord.com/api/v9/applications/{client_id}/commands'
    with open(json_path, 'r') as file:
        commands = file.read()
        response = requests.put(url, headers=headers, data=commands)
        try:
            print(response.json())
        except:
            print(response.text)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Register commands to Discord API')
    parser.add_argument('--json', '-j', type=str, help='Path to the JSON file containing the commands', required=True)
    parser.add_argument('--bot-token', '-t', type=str, help='Bot token for authentication', required=True)
    parser.add_argument('--client-id', '-i', type=str, help='Client ID for the application', required=True)
    
    args = parser.parse_args()

    register_commands(args.json, args.bot_token, args.client_id)