# Getting Started with Sysl

This is a basic guide to building and using Sysl. It is intended for early adopters who want to tinker with Sysl, but need a little guidance getting started.


## Installation

Currently the best way to start using Sysl, is to build it yourself from source. 

### Windows
TODO : Windows build instructions

### Linux
TODO : Linux build instructions

### OSX

Begin by cloning the repo to your local machine
```bash
git clone https://github.com/anz-bank/sysl.git
cd sysl
```

You will need to install [Bazel](https://www.bazel.build/) in order to build the compiler. The simplest way is via [homebrew](https://brew.sh/)
```bash
brew install bazel
```

Next you will need to install the required python dependencies.
```bash
sudo easy_install pip
sudo -H pip install --ignore-installed six
sudo -H pip install plantuml requests protobuf
```

Finally, you can build the compiler using Bazel
```bash
bazel build //src/sysl
```

By default, Bazel will output to ```bazel-bin/src/sysl/``` 

To test that everything worked, you can output the default ```sysl``` help message by running:
```
bazel-bin/src/sysl/sysl -h
```

To add ```sysl``` to your $PATH, run the following:

```bash
printf "\nexport PATH=$(PWD)/bazel-bin/src/sysl:\$PATH" >> ~/.bash_profile
source ~/.bash_profile
```

You should now be able to run ```sysl -h``` from anywhere on your machine.

## Using the Sysl compiler

At this point, you should have successfully built ```sysl``` and added it to your path. Now comes the fun part...

In the root of the project you will find a folder named ```demo``` containing various examples of Sysl modules.  For this guide, we will be focusing on the ```demo/petshop``` example. 

If you haven't already, open up ```demo/petshop/petshop.sysl``` and have a quick look at the module we will be compiling. 

Before we begin, let's create a temporary folder to hold our output files
```
mkdir gsguide
```

### Compiling your first module

To build compile the ```demo/petshop``` module, run the following command
```bash
sysl textpb --output gsguide/petshop.txt //demo/petshop/petshop
```

Congratulations! You just compiled your very first Sysl module.

If you open ```gsguide/petshop.txt``` with your favorite editor, you will see the result of all your hard work.

Let's examine each part of that command

|                                  |                                                        |
|----------------------------------|--------------------------------------------------------|
|```sysl textpb```                 |Instructs ```sysl``` to use the ```textpb``` sub-command|
|```--output gsguide/petshop.txt```|Specifies the path of the output file                               |
|```//demo/petshop/petshop```      |Which module ```sysl``` will compile                              |


### Sub-commands

When parsing command line arguments to ```sysl```, the first argument must always be a valid sub-command.

Each sub-command will produce a different output format, and some sub-commands require specific command line arguments.

Going over the sub-commands in detail is beyond the scope of this guide, but here is a quick summary

| sub-command |  output  |
|-------------|----------|
|pb           | [Protocol buffers](https://github.com/google/protobuf)|
|textpb       | Text serialized [Protocol buffers](https://github.com/google/protobuf)|
|deps         | 
|data         | 
|ints         | 
|sd           |

TODO : Document the output format, and required command line arguments for each sub-command

### Modules

The final argument we passed to ```sysl``` was the path of our module  
```//demo/petshop/petshop```

For our example we only needed to compile a single module, but in actuality ```sysl``` can compile multiple modules at once, provided they are all listed at the **end** of the argument list.

Modules can be referenced as follows :

- Module names always start with ```//```
- Names are based on the file path
- The ```.sysl``` file extension is automatically appended

Applying this to our example,  ```demo/petshop/petshop.sysl``` becomes ```//demo/petshop/petshop```

## Conclusion

We've barely scratched the surface of what Sysl can do!

Hopefully this guide has given you a decent understanding of the basics of compiling with ```sysl```.

Unfortunately, Sysl is still in early alpha, and lacks solid documentation. This will improve as the project matures out of alpha, but for now at least you are on your own.

If you have questions or need help, you can try asking [here](https://github.com/anz-bank/sysl/issues) 

Good luck!

TODO: Links to more resources / documentation