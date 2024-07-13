package main

import (
    "flag"
    "encoding/csv"
    "fmt"
    "os"
    "time"
)

func problemPuller(fileName string)([]problem, error) {
    if fObj, err := os.Open(fileName); err == nil {
        csvReader := csv.NewReader(fObj) 
        if cLines , err := csvReader.ReadAll(); err == nil {
            return problemParser(cLines), nil
        }else {
            return nil , fmt.Errorf("unable to read the file" + "file from %s file; %s", fileName, err.Error())
        }
    }else {
        return nil, fmt.Errorf("error opening the: %s file; %s", fileName, err.Error())
    }
}

func main() {
    fName := flag.String("f", "quiz.csv", "path of the file")

    timer := flag.Int("t", 30, "timer of the quiz")
    flag.Parse()

    problems, err := problemPuller(*fName)

    if err != nil {
        exit(fmt.Sprintf("an error occured:%s", err.Error()))
    }

    correctAnswer := 0

    tObjs := time.NewTimer(time.Duration(*timer) * time.Second) 
    ansC := make(chan string)

    problemLoop:

     for i, p := range problems{
        var answer string
        fmt.Printf("Problem %d: %s=", i+1, p.q)

        go func ()  {
            fmt.Scanf("%s", &answer)
            ansC <- answer
        }()
        select{
        case <- tObjs.C:
        fmt.Println()
        break problemLoop
        case iAns := <-ansC:
        if iAns == p.a{
                correctAnswer++
            }
        if i == len(problems)-1{
close(ansC)
            }
    }
    }
    fmt.Printf("Your result is %d out of %d\n", correctAnswer, len(problems))
    fmt.Printf("Press enter to exit")
}

func problemParser(lines [][]string) []problem {

    r := make([]problem, len(lines))

    for i := 0 ; i < len(lines); i++ {
        r[i] = problem{q: lines[i][0], a: lines[i][1]}
    }
    return r
}

type problem struct {
    q string
    a string
}
func exit(msg string) {
    fmt.Println(msg)
    os.Exit(1)
}
