package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type Question struct {
	ID           int
	Question     string
	Options      [4]string
	Answer       string
	CorrectCount int
	WrongCount   int
}

type Participant struct {
	ID    int
	Name  string
	Score int
}

var (
	questions     []Question
	participants  []Participant
	currentID     int = 1
	currentQuesID int = 1
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		fmt.Println("Welcome to the Quiz Game!")
		fmt.Println("1. Admin Login")
		fmt.Println("2. Participant Login")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			adminMenu()
		case 2:
			participantMenu()
		case 3:
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func adminMenu() {
	for {
		fmt.Println("\n--- Admin Menu ---")
		fmt.Println("1. Add Question")
		fmt.Println("2. Edit Question")
		fmt.Println("3. Delete Question")
		fmt.Println("4. View Top 5 Most Correct/Wrong Questions")
		fmt.Println("5. Back")
		fmt.Print("Choose an option: ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			addQuestion()
		case 2:
			editQuestion()
		case 3:
			deleteQuestion()
		case 4:
			viewTopQuestions()
		case 5:
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func addQuestion() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter question: ")
	question, _ := reader.ReadString('\n')
	question = strings.TrimSpace(question)

	var options [4]string
	for i := 0; i < 4; i++ {
		fmt.Printf("Enter option %c: ", 'A'+i)
		option, _ := reader.ReadString('\n')
		options[i] = strings.TrimSpace(option)
	}

	fmt.Print("Enter the correct answer (A/B/C/D): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	questionData := Question{
		ID:       currentQuesID,
		Question: question,
		Options:  options,
		Answer:   answer,
	}
	questions = append(questions, questionData)
	currentQuesID++

	fmt.Println("Question added successfully!")
}

func editQuestion() {
	var id int
	fmt.Print("Enter question ID to edit: ")
	fmt.Scanln(&id)

	for i, q := range questions {
		if q.ID == id {
			var question, answer string
			var options [4]string
			fmt.Print("Enter new question: ")
			fmt.Scanln(&question)
			fmt.Print("Enter new option A: ")
			fmt.Scanln(&options[0])
			fmt.Print("Enter new option B: ")
			fmt.Scanln(&options[1])
			fmt.Print("Enter new option C: ")
			fmt.Scanln(&options[2])
			fmt.Print("Enter new option D: ")
			fmt.Scanln(&options[3])
			fmt.Print("Enter new correct answer (A/B/C/D): ")
			fmt.Scanln(&answer)

			questions[i] = Question{
				ID:       id,
				Question: question,
				Options:  options,
				Answer:   answer,
			}
			fmt.Println("Question updated successfully!")
			return
		}
	}
	fmt.Println("Question not found!")
}

func deleteQuestion() {
	var id int
	fmt.Print("Enter question ID to delete: ")
	fmt.Scanln(&id)

	for i, q := range questions {
		if q.ID == id {
			questions = append(questions[:i], questions[i+1:]...)
			fmt.Println("Question deleted successfully!")
			return
		}
	}
	fmt.Println("Question not found!")
}

func viewTopQuestions() {
	fmt.Println("\n--- Top 5 Most Correct/Wrong Questions ---")
	fmt.Println("1. View Most Correct Answers")
	fmt.Println("2. View Most Incorrect Answers")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scanln(&choice)

	var sortedQuestions []Question
	switch choice {
	case 1:
		sortedQuestions = make([]Question, len(questions))
		copy(sortedQuestions, questions)
		sort.Slice(sortedQuestions, func(i, j int) bool {
			return sortedQuestions[i].CorrectCount > sortedQuestions[j].CorrectCount
		})
	case 2:
		sortedQuestions = make([]Question, len(questions))
		copy(sortedQuestions, questions)
		sort.Slice(sortedQuestions, func(i, j int) bool {
			return sortedQuestions[i].WrongCount > sortedQuestions[j].WrongCount
		})
	default:
		fmt.Println("Invalid choice, going back.")
		return
	}

	for i := 0; i < 5 && i < len(sortedQuestions); i++ {
		fmt.Printf("ID: %d, Question: %s\n", sortedQuestions[i].ID, sortedQuestions[i].Question)
	}
}

func participantMenu() {
	fmt.Print("Enter your name: ")
	var name string
	fmt.Scanln(&name)

	participant := Participant{
		ID:   currentID,
		Name: name,
	}
	currentID++

	participants = append(participants, participant)

	playQuiz(participant)
}

func playQuiz(participant Participant) {
	fmt.Println("\n--- Starting Quiz ---")
	score := 0
	questionCount := 5

	shuffledQuestions := shuffleQuestions()

	for i := 0; i < questionCount && i < len(shuffledQuestions); i++ {
		q := shuffledQuestions[i]
		fmt.Printf("\nQ%d: %s\n", i+1, q.Question)
		for j, option := range q.Options {
			fmt.Printf("%c. %s\n", 'A'+j, option)
		}

		var answer string
		fmt.Print("Your answer (A/B/C/D): ")
		fmt.Scanln(&answer)

		if strings.ToUpper(answer) == q.Answer {
			score++
			questions[i].CorrectCount++
		} else {
			questions[i].WrongCount++
		}

		fmt.Printf("Correct Answer: %s\n", q.Answer)
	}

	participant.Score = score
	fmt.Printf("\nYour final score is: %d/%d\n", score, questionCount)

	updateParticipantScore(participant)

	displayLeaderboard()
}

func shuffleQuestions() []Question {
	shuffled := make([]Question, len(questions))
	copy(shuffled, questions)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func updateParticipantScore(participant Participant) {
	for i, p := range participants {
		if p.ID == participant.ID {
			participants[i].Score = participant.Score
			return
		}
	}
}

func displayLeaderboard() {
	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Score > participants[j].Score
	})

	fmt.Println("\n--- Leaderboard ---")
	for i, p := range participants {
		fmt.Printf("%d. %s - Score: %d\n", i+1, p.Name, p.Score)
	}
}
