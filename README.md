# Couples Quiz App - Backend Service

## Overview

This project is the backend service powering a Couples Quiz application. It provides the necessary API endpoints and core logic to manage quiz content, including categories, quizzes, and various types of interactive questions. The system is designed to handle quiz creation, updates, retrieval, and associated media like images.

## Core Features

*   **Quiz Category Management:** Administrators can create, update, and delete quiz categories, including uploading category-specific images.
*   **Quiz Management:** Users (likely authenticated) can create, update, retrieve, and delete quizzes. This includes:
    *   Associating quizzes with categories.
    *   Adding descriptive text and images.
    *   Support for different languages.
    *   Publishing quizzes to make them available.
*   **Diverse Question Types:** Supports the creation and management of various question formats within a quiz, such as:
    *   True/False
    *   Slider
    *   Ordering
    *   Open Answer
    *   Multiple Choice (Single/Multiple Answer)
    *   Matching
    *   Drag and Drop
*   **Image Handling:** Integrates with a file service to upload, manage, and retrieve images associated with categories, quizzes, and even specific question options.
*   **Data Retrieval:** Offers flexible ways to fetch quizzes and categories, including filtering and pagination.
*   **Authorization:** Includes checks to ensure only authorized users (e.g., the quiz creator) can modify specific quizzes or questions.

## How It Works (High-Level)

The backend is built using **Go** and follows a layered architecture pattern:

1.  **Service Layer (`appquizzes`):** Contains the core business logic. It orchestrates operations by interacting with repositories and other services (like file handling or user services). It defines distinct services for Admin (`AdminService`) and User (`UserService`) functionalities.
2.  **Repository Layer (`repository.go`):** Defines interfaces for data persistence operations (CRUD for quizzes, categories, questions, etc.). This abstracts the underlying database interactions.
3.  **Domain Layer (`domain.go`, `models.go`, `errors.go`):** Defines the core entities (structs like `QuizModel`, `QuestionPlainModel`), service interfaces, custom errors, constants (like question types), and request/response structures.
4.  **Infrastructure/Dependencies:**
    *   Uses interfaces for dependencies like `Transaction` management (ensuring atomicity for complex operations like deletions), `FileService` (for image uploads/URLs), and `LocalizationService`.
    *   Handles different question option formats using JSON and specific creator/deletor logic for each question type.
    *   Employs dependency injection to wire components together (evident in `NewUserService` and `NewAdminServiceImpl`).

This structure promotes separation of concerns, testability, and maintainability. The use of specific handlers for different question types allows for easy extension with new formats in the future.
