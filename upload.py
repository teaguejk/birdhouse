# imports
from fabric.connection import Connection

# setup
user = 'teaguejk'
host = 'student2.cs.appstate.edu'
path = '/usr/local/apache2/htdocs/u/teaguejk/birdhouse.site'
filename = './assets/IMG.jpg'

# main
def main():
    print("For server: " + user + "@" + host)
    spass = input("Server Password:\n")

    with Connection(host, user, connect_kwargs={'password': spass, 'allow_agent': False}) as c:  
        c.put(filename, path + "/assets/IMG.jpg")

if __name__ == "__main__":
    main()

