
from fabric.connection import Connection
user = 'teaguejk'
host = 'student2.cs.appstate.edu'
path = '/usr/local/apache2/htdocs/u/teaguejk/birdhouse.site'
filename = './assets/IMG.jpg'
with Connection(host, user, connect_kwargs={'password': '900728429', 'allow_agent': False}) as c:  
         c.put(filename, path + "/assets/IMG.jpg")