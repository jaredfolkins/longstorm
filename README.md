# 🌩️ LongStorm 🌩️

LongStorm is a web application that you run locally and that helps faciliate the posting of long form threaded tweets (tweet storms).

# Releases

- Download the latest release on the [releases page](https://github.com/jaredfolkins/longstorm/releases).
- Set the executable bit ```$ chmod 744 longstorm```
- Run it from the cmd line ```$ ./longstorm```
- Open up a web browser and go to http://localhost:3000

# Compile

```
$ export GO111MODULE=on
$ git clone https://github.com/jaredfolkins/longstorm
$ cd longstorm
$ go run main.go
```

Now open your web browser to http://localhost:3000

# Screenshots

### take a wall of text
![longstorm wall of text](https://raw.githubusercontent.com/jaredfolkins/longstorm/master/assets/images/wall_of_text.png)

### And you can create a threaded twitter post

![longstorm wall of text](https://raw.githubusercontent.com/jaredfolkins/longstorm/master/assets/images/longstorm_threaded_app.png)

# Twitter API Keys 

You must use your account to create twitter api developer keys and save those keys into LongStorm. Be aware the keys are sensitive so don't lose the them or display them!

### Generate the Twitter API Keys

![longstorm wall of text](https://raw.githubusercontent.com/jaredfolkins/longstorm/master/assets/images/twitter_keygen_example.png)

### Save the keys to LongStorm

![longstorm wall of text](https://raw.githubusercontent.com/jaredfolkins/longstorm/master/assets/images/input_twitter_keys_example.png)
