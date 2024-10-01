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

	# for testing purposes
	# channel = client.get_channel(796549061200314432)
	# await channel.send('hi')

	game = discord.Game("watching for birds")
	await client.change_presence(status=discord.Status.idle, activity=game)


@client.event
async def on_message(message):

	if message.author == client.user:
		return

	if message.content == '!bird':
		# get the most recent image
		files = glob.glob(path + filetype)
		recent = max(files, key=os.path.getctime)
		with open(recent, 'rb') as f:
			img = discord.File(f)
			await message.channel.send(file=img)
	else if message.content == '!help':
		await message.channel.send("Use !bird to display the most recent image...")

client.run(token)
