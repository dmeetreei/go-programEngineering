package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Структура данных сотрудника
type Employee struct {
	Fullname   string  `json:"fullname"`
	Position   string  `json:"position"`
	Experience float64 `json:"experience"`
}

type Employees []Employee

// Добавить сотрудника
func (employees *Employees) addEmployee(fullname, position string, experience float64) {
	*employees = append(*employees, Employee{
		Fullname:   fullname,
		Position:   position,
		Experience: experience,
	})
}

// Удалить сотрудника
func (employees *Employees) removeEmployee(fullname string) {
	for i, emp := range *employees {
		if emp.Fullname == fullname {
			fmt.Printf("Удаляется сотрудник %d", i+1)
			*employees = append((*employees)[:i], (*employees)[i+1:]...)
			return
		}
	}
}

// Показать сотрудников
func displayEmployees(employees *[]Employee) {
	fmt.Printf("%-5s | %-30s | %-20s | %-10s\n", "№", "Ф.И.О.", "Должность", "Стаж (годы)")
	fmt.Println(strings.Repeat("-", 70))
	for i, employee := range *employees {
		fmt.Printf("%-5d | %-30s | %-20s | %-10.2f\n", i+1, employee.Fullname, employee.Position, employee.Experience)
	}
}

// Сортирует массив сотрудников по алфавиту в поле Ф.И.О.
func sortEmployees(employees *Employees) {
	sort.SliceStable(*employees, func(i, j int) bool {
		return strings.ToLower((*employees)[i].Fullname) < strings.ToLower((*employees)[j].Fullname)
	})
	fmt.Println("Сохранение файла employees.json")
	err := SaveToFile(*employees)
	if err != nil {
		return
	}
}

// Сортирует массив сотрудников по алфавиту в поле Должность
func sortPositions(employees *Employees) {
	sort.SliceStable(*employees, func(i, j int) bool {
		return strings.ToLower((*employees)[i].Position) < strings.ToLower((*employees)[j].Position)
	})
	fmt.Println("Список сотрудников отсортирован по должностям.")
	displayEmployees((*[]Employee)(employees))
}

// Вычисляет средний стаж и выводит сотрудников, работающих на должности с минимальным средним стажем.
// kostyl
func findPosition(employees Employees) {
	averageExperience := calculatePosition(employees)

	if len(averageExperience) == 0 {
		fmt.Println("Нет сотрудников для анализа.")
		return
	}

	var minPosition string
	var minExperience float64 = -1

	// Найдем должность с минимальным средним стажем
	for position, experience := range averageExperience {
		if minExperience == -1 || experience < minExperience {
			minExperience = experience
			minPosition = position
		}
	}

	fmt.Printf("Должность с минимальным средним стажем: %s (Средний стаж: %.2f)\n", minPosition, minExperience)
	fmt.Println("Сотрудники на этой должности:")

	// Выведем всех сотрудников, работающих на этой должности
	for _, emp := range employees {
		if emp.Position == minPosition {
			fmt.Printf("Ф.И.О.: %s, Должность: %s, Стаж работы: %.2f\n", emp.Fullname, emp.Position, emp.Experience)
		}
	}
}

// kostyl
func calculatePosition(employees Employees) map[string]float64 {
	positionExperience := make(map[string]float64)
	positionCount := make(map[string]int)

	for _, emp := range employees {
		positionExperience[emp.Position] += emp.Experience
		positionCount[emp.Position]++
	}

	averageExperience := make(map[string]float64)
	for position, totalExperience := range positionExperience {
		averageExperience[position] = totalExperience / float64(positionCount[position])
	}

	return averageExperience
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func SaveToFile(employees []Employee) error {
	file, err := os.Create("employees.json")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Ошибка при сохранении файла")
		}
	}(file)

	encoder := json.NewEncoder(file)
	fmt.Println("Файл сохранен")
	return encoder.Encode(employees)
}

func LoadFromFile(filename string) ([]Employee, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var employees Employees
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&employees); err != nil {
		return nil, fmt.Errorf("ошибка декодирования данных: %v", err)
	}
	return employees, nil
}

func main() {
	const filename = "employees.json"
	var employees Employees

	if FileExists(filename) {
		fmt.Println("Найден файл employees.json. Хотите загрузить данные из него или очистить его?")
		fmt.Println("Введите '1' для загрузки данных или '2' для очистки файла:")
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil || (choice != 1 && choice != 2) {
			fmt.Println("Некорректный ввод, программа завершена.")
			return
		}

		if choice == 1 {
			// Загрузка данных из файла
			loadedEmployees, err := LoadFromFile(filename)
			if err != nil {
				fmt.Printf("Ошибка при загрузке данных: %v\n", err)
				return
			}
			employees = loadedEmployees
			fmt.Println("Данные успешно загружены.")
		} else if choice == 2 {
			// Очистка файла (создание нового пустого файла)
			err := os.Remove(filename)
			if err != nil {
				fmt.Printf("Ошибка при очистке файла: %v\n", err)
				return
			}
			fmt.Println("Файл успешно очищен. Начинаем работу с пустым списком сотрудников.")
		}
	} else {
		// Файл не существует, создаем новый
		fmt.Println("Файл employees.json не найден. Создаем новый и начинаем работу.")
	}

	var a = 1
	for a != 0 {
		fmt.Println("Меню")
		fmt.Println("1. Добавить сотрудника")
		fmt.Println("2. Удалить сотрудника")
		fmt.Println("3. Отсортировать сотрудников по алфавиту")
		fmt.Println("4. Отсортировать сотрудников по должности")
		fmt.Println("5. Показать сотрудников на должности с минимальным стажем")
		fmt.Println("6. Показать средний стаж по каждой должности")
		fmt.Println("9. Показать всех сотрудников")
		fmt.Println("0. Сохранить и выйти")
		_, err := fmt.Scanln(&a)
		if err != nil {
			return
		}
		switch a {
		case 1:
			reader := bufio.NewReader(os.Stdin)

			var fullname, position string
			var experience float64

			fmt.Println("Введите Ф.И.О.")
			fullname, _ = reader.ReadString('\n')
			fullname = strings.TrimSpace(fullname)

			fmt.Println("Введите должность")
			position, _ = reader.ReadString('\n')
			position = strings.TrimSpace(position)

			fmt.Println("Введите стаж работы (в годах, дробное число):")
			_, err := fmt.Scanln(&experience)
			if err != nil {
				return
			}
			employees.addEmployee(fullname, position, experience)
		case 2:
			reader := bufio.NewReader(os.Stdin)
			var fullname string

			fmt.Println("Введите Ф.И.О.")
			fullname, _ = reader.ReadString('\n')
			fullname = strings.TrimSpace(fullname)

			employees.removeEmployee(fullname)
		case 3:
			fmt.Println("Сортировка по алфавиту...")
			sortEmployees(&employees)
		case 4:
			sortPositions(&employees)
		case 5:
			findPosition(employees)

		case 6:
			averageExperience := calculatePosition(employees)
			fmt.Println("Средний стаж по каждой должности:")
			for position, experience := range averageExperience {
				fmt.Printf("Должность: %s, Средний стаж: %.2f\n", position, experience)
			}
		case 9:
			displayEmployees((*[]Employee)(&employees))
		case 0:
			fmt.Println("Сохранение файла employees.json")
			err := SaveToFile(employees)
			if err != nil {
				return
			}
			fmt.Println("Выход...")
			os.Exit(0)
		}
	}
}
