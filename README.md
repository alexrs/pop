# Pop
StackOverflow answers in your terminal.

## Usage

```TEXT
$ pop open file java

Read file, parse each line into an integer and store into a list:

List<Integer> list = new ArrayList<Integer>();
File file = new File("file.txt");
BufferedReader reader = null;

try {
    reader = new BufferedReader(new FileReader(file));
    String text = null;

    while ((text = reader.readLine()) != null) {
        list.add(Integer.parseInt(text));
    }
} catch (FileNotFoundException e) {
    e.printStackTrace();
} catch (IOException e) {
    e.printStackTrace();
} finally {
    try {
        if (reader != null) {
            reader.close();
        }
    } catch (IOException e) {
    }
}

//print out the list
System.out.println(list);


Url: http://stackoverflow.com/questions/3806062/how-to-open-a-txt-file-and-read-numbers-in-java
```

## Install

```BASH
go get github.com/alexrs/pop
```

## TODO
- [ ] Parse the HTML better and display it properly
- [ ] Tests
- [ ] More options (contribute or suggest new features!)
