package handlers

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"user/app/rpc/pb"
	. "user/app/rpc/pb"
	"user/core/dependencies/services"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

//NewUserHandler factory method for creating a method handler for users
func NewUserHandler(depedencies services.App) AppHandler {
	return userHandler{appDependencies: depedencies}
}

const (
	imagePath = "/api/v1/users/image/"
)

type userHandler struct {
	AppHandler
	pb.UnsafeUserServiceServer
	appDependencies services.App
}

// func (object userHandler) userImage() services.ProfileImageServices {
// 	return object.appDependencies.Image.UserProfileImage
// }

// func (object userHandler) userValidator() services.UserValidatorServices {
// 	return object.appDependencies.Validator.UserValidator
// }

// func (object userHandler) appCredentials() config.Credentials {
// 	return object.appDependencies.Credentials
// }

// func (object userHandler) userRepository() entities.UserRepository {
// 	return object.appDependencies.Repository.UserRepository
// }

// type BasicUser struct {
// 	Name  string
// 	Email string
// }

func (object *userHandler) RegisterCrud(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, object)
	// app.Get("/api/v1/users", object.getUsers)
	// app.Get("/api/v1/users/email", object.getUser)
	// app.Get(imagePath+":id", object.userImage().NewDownloader(), object.getImage)
	// app.Post("/api/v1/users", object.userValidator().DuplicatedUser(), object.userImage().NewUploader(), object.newUser)
	// app.Put("/api/v1/users", object.userImage().UpdateImage(), object.updateUser)
	// app.Delete("/api/v1/users/email", object.userImage().DeleteImage(), object.deleteUser)
}

func (object *userHandler) GetUser(context.Context, *UserIdentity) (*UserResponse, error) {
	response := &UserResponse{
		Body: &User{
			Name:        "example",
			Email:       "tes@example.com",
			Nickname:    "test 1",
			Password:    "",
			ImagePath:   "",
			CountryCode: "COL",
			Birthday:    "03/09/1988",
		},
	}

	return response, nil
}

func (object *userHandler) GetUsers(context.Context, *emptypb.Empty) (*UsersResponse, error) {
	allUsers := []*User{
		&User{
			Name:        "example",
			Email:       "tes@example.com",
			Nickname:    "test 1",
			Password:    "",
			ImagePath:   "",
			CountryCode: "COL",
			Birthday:    "03/09/1988",
		}, &User{
			Name:        "example2",
			Email:       "tes@example2.com",
			Nickname:    "test 2",
			Password:    "",
			ImagePath:   "",
			CountryCode: "COL",
			Birthday:    "03/09/1988",
		},
	}

	response := &UsersResponse{
		Users: allUsers,
	}

	return response, nil
}

func (object *userHandler) SaveUser(context.Context, *UserRequest) (*UserResponse, error) {
	response := &UserResponse{
		Body: &User{
			Name:        "example",
			Email:       "tes@example.com",
			Nickname:    "test 1",
			Password:    "",
			ImagePath:   "",
			CountryCode: "COL",
			Birthday:    "03/09/1988",
		},
	}

	return response, nil
}

func (object *userHandler) UpdateUser(context.Context, *UserRequest) (*UserResponse, error) {
	response := &UserResponse{
		Body: &User{
			Name:        "example",
			Email:       "tes@example.com",
			Nickname:    "test 1",
			Password:    "",
			ImagePath:   "",
			CountryCode: "COL",
			Birthday:    "03/09/1988",
		},
	}

	return response, nil
}

func (object *userHandler) SaveProfileImage(fileStreamRequest UserService_SaveProfileImageServer) error {
	buffer, filename, err := processStream(fileStreamRequest)
	if err != nil {
		log.Printf("Error receiving chunk data, %v", err)
		return err
	}

	err = processBuffer(buffer, filename)
	if err != nil {
		log.Printf("Error creating file, %v", err)
		return err
	}

	response := &UserResponse{
		Body: &User{
			Name:        "example",
			Email:       "tes@example.com",
			Nickname:    "test 1",
			Password:    "",
			ImagePath:   "",
			CountryCode: "COL",
			Birthday:    "03/09/1988",
		},
	}

	fileStreamRequest.SendAndClose(response)
	return nil
}

func processStream(fileStreamRequest pb.UserService_SaveProfileImageServer) (*bytes.Buffer, string, error) {
	buffer := &bytes.Buffer{}
	filename := ""

	for {
		chunkRequest, err := fileStreamRequest.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, "", err
		}

		log.Printf("Receive new chunk data, secuence %d from total %d", chunkRequest.CurrentSecuence, chunkRequest.LastSecuence)

		if _, err = buffer.Write(chunkRequest.Content); err != nil {
			return nil, "", err
		}

		if filename == "" {
			filename = chunkRequest.Filename
		}
	}

	return buffer, filename, nil
}

//TODO: Send this buffer(new image file) to S3 helper
func processBuffer(buffer *bytes.Buffer, videoName string) error {
	basePath := "../../assets/server/"
	videoFile := basePath + videoName
	os.Remove(videoFile)

	if err := ioutil.WriteFile(videoFile, buffer.Bytes(), 777); err != nil {
		return err
	}

	log.Println("File " + videoName + " created successfully")
	return nil
}

// func (object userHandler) getUsers(context *fiber.Ctx) error {
// 	users, err := object.userRepository().Users()

// 	if err != nil {
// 		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
// 	}

// 	if len(users) == 0 {
// 		return response.MakeSuccessJSON([]entities.User{}, context)
// 	}

// 	return response.MakeSuccessJSON(users, context)
// }

// func (object userHandler) getUser(context *fiber.Ctx) error {
// 	userEmail := context.Query("address")

// 	if userEmail == "" {
// 		return response.MakeErrorJSON(http.StatusBadRequest, "address is not present on url as a query param")
// 	}

// 	user, err := object.userRepository().User(userEmail)

// 	if err != nil {
// 		return response.MakeErrorJSON(http.StatusNotFound, err.Error())
// 	}

// 	return response.MakeSuccessJSON(user, context)
// }

// func (object userHandler) getImage(context *fiber.Ctx) error {
// 	if s3DataFile, ok := context.Locals(image.PROFILE_IMAGE_DOWNLOADED_FILE).(*services.ImageBufferedFile); ok {
// 		return context.Status(http.StatusOK).SendStream(bytes.NewReader(s3DataFile.Data), int(s3DataFile.Size))
// 	}

// 	return response.MakeErrorJSON(http.StatusInternalServerError, "Invalid type file")
// }

// func (object userHandler) newUser(context *fiber.Ctx) error {
// 	user, err := getUserRequestBody(object, context)
// 	if err != nil {
// 		return err
// 	}

// 	if err := object.userRepository().CreateUser(user); err != nil {
// 		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
// 	}

// 	return response.MakeSuccessJSON("user was created successfully", context)
// }

// func (object userHandler) updateUser(context *fiber.Ctx) error {
// 	user, err := getUserRequestBody(object, context)
// 	if err != nil {
// 		return err
// 	}

// 	if err := object.userRepository().UpdateUser(user); err != nil {
// 		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return response.MakeSuccessJSON("user updated successfully", context)
// }

// func (object userHandler) deleteUser(context *fiber.Ctx) error {
// 	var user entities.User

// 	if assertion, ok := context.Locals(image.PROFILE_IMAGE_USER_ENTITY).(entities.User); ok {
// 		user = assertion
// 	}

// 	if err := object.userRepository().DeleteUser(user.Email); err != nil {
// 		return response.MakeErrorJSON(http.StatusInternalServerError, "Invalid user")
// 	}

// 	return response.MakeSuccessJSON("user deleted successfully", context)
// }

// func getUserRequestBody(object userHandler, context *fiber.Ctx) (*entities.User, error) {
// 	user := new(entities.User)

// 	if err := context.BodyParser(user); err != nil {
// 		return nil, response.MakeErrorJSON(http.StatusBadRequest, err.Error())
// 	}

// 	if imageURI, ok := context.Locals(image.PROFILE_IMAGE__UPLOADED_ID).(string); ok {
// 		user.ImageID = imagePath + imageURI
// 	} else {
// 		user.ImageID = imagePath + image.DEFAULT_IMAGE
// 	}

// 	if err := object.userValidator().IsValid(*user); err != nil {
// 		return nil, response.MakeErrorJSON(http.StatusBadRequest, err.Error())
// 	}

// 	return user, nil
// }
