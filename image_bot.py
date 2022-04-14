"""
Jaracah Teague
image_bot.py

Discord Bot for my birdhouse project
Sends the most recent image to a discord channel

"""

import discord
import os.path
import glob

client = discord.Client()
token = ''
path = './assets/'
filetype = file_type = r'/*jpg'

@client.event
async def on_ready():
	print('Logged in as {0.user}'.format(client))
	channel = client.get_channel(796549061200314432)
	# await channel.send('hi')

	files = glob.glob(path + filetype)
	recent = max(files, key=os.path.getctime)
	with open(recent, 'rb') as f:
	    img = discord.File(f)
	    await channel.send(file=img)

client.run(token)

