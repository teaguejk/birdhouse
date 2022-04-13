"""
Jaracah Teague
image_bot.py

Discord Bot for my birdhouse project
Sends the most recent image to a discord channel

"""

import discord

client = discord.Client()
token = ''

@client.event
async def on_ready():
	print('Logged in as {0.user}'.format(client))
	channel = client.get_channel(796549061200314432)
	await channel.send('hi')

client.run(token)

