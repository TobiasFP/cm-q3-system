# cm-q3-system

This is my take on the question 3 "Design a system to update robot's maps in stores"

## Notes

- Many things have been hardcoded, such as mappings to ports, etc. These should ideally by written in a config, but as this is only supposed to be pseudo code, I have in many places hardcoded theses.
- I have only set up a HTTP/1.1 server, not HTTP/2, as I am unaware of any feature within http/2 that could make this software any better.
- My unittests are crude, as I had very little time to make this setup. Tests will fail if the code has been run. In order to have running tests,
  you need a freshly seeded DB.
- The Go Code is not neatly put into packages, as I would normally do. This is again due to the time constraints, and the size of the project.
- I am aware that the map that I have used is a picture of a map. This is not because I do not understand the concept of binary maps, it is purely because
  i do not have such binary maps readily available, and the picture seamed appropriate to work with.
- The docker compose file has networks to simulate the isolation of each system. My idea was to have a master container that could access all
  networks separately. I have not had the time to implement this, and I believe using the ports to access the various parts is sufficient to illustrate the isolation of each system.

## Setting up

- docker-compose up -d
- cd src
- ./cm-q3-watcher_linux || ./cm-q3-watcher_mac || cm-q3-watcher.exe
