package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/zrwaite/Insomnizac/db"
	"github.com/zrwaite/Insomnizac/graph/model"
	"github.com/zrwaite/Insomnizac/graph/services/queries"
	"github.com/zrwaite/Insomnizac/utils"
	"github.com/zrwaite/Insomnizac/utils/httpreq"
)

var defaultImage = "https://storage.googleapis.com/insomnizac_public/static/default_project.png"
var emptyMap = map[string]string{}

func GetProjectArgs(project *model.Project) []interface{} {
	return []interface{}{&project.ID, &project.Name, &project.Slug, &project.GithubName, &project.DevpostLink, &project.ProjectLink, &project.CreatedAt, &project.UpdatedAt, &project.Image, &project.Featured, pq.Array(&project.SkillIds)}
}

func GetGithubProject(project *model.Project) {

}

func GetProject(slug string) (project *model.Project, status int) {
	cacheKey := "project:" + slug
	project = new(model.Project)
	found := db.GetJsonCache(cacheKey, project)
	if found {
		return project, 200
	}

	row := db.DB.QueryRow("SELECT * FROM projects WHERE slug=$1", slug)
	if row.Err() != nil {
		return nil, 400
	}
	err := row.Scan(GetProjectArgs(project)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 404
		} else {
			log.Fatal(err)
		}
	}
	if project.Image == nil {
		project.Image = &defaultImage
	}
	project.ParseProject()
	err = GetRepositoryData(project)
	if err != nil {
		fmt.Println(err)
		return nil, 400
	}
	db.SetJsonCache(cacheKey, project)
	return project, 200
}

func GetProjects() (projects []*model.Project, status int) {
	cacheKey := "projects"
	found := db.GetJsonCache(cacheKey, &projects)
	if found {
		return projects, 200
	}

	rows, err := db.DB.Query("SELECT * FROM projects")
	if err != nil {
		return nil, 400
	}
	defer rows.Close()
	for rows.Next() {
		project := new(model.Project)
		err = rows.Scan(GetProjectArgs(project)...)
		if err != nil {
			fmt.Println(err)
			return nil, 400
		}
		if project.Image == nil {
			project.Image = &defaultImage
		}
		project.ParseProject()
		projects = append(projects, project)
	}
	err = GetRepositoriesData(projects)
	if err != nil {
		log.Fatal(err)
	}
	db.SetJsonCache(cacheKey, projects)
	return projects, 200
}

func GetRepositoryData(project *model.Project) error {
	username, projectName := utils.GetProjectNames(project.GithubName)
	variables := map[string]string{"owner": username, "name": projectName}
	resp, err := httpreq.GithubQuery(queries.RepositoryQuery, variables)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("github API returned non-200 status code")
	}
	var body model.GithubRepoResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}
	project.Description = body.Data.Repository.Description
	return nil
}

func GetRepositoriesData(projects []*model.Project) error {
	repoGithubNames := []string{}
	for _, project := range projects {
		repoGithubNames = append(repoGithubNames, project.GithubName)
	}
	repoQuery := queries.GenerateRepositoriesQuery(repoGithubNames)
	resp, err := httpreq.GithubQuery(repoQuery, emptyMap)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("github API returned non-200 status code")
	}
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}
	data := body["data"].(map[string]interface{})
	for _, project := range projects {
		key := utils.GetQueryName(project.GithubName)
		repo := data[key].(map[string]interface{})
		project.Description = repo["description"].(string)
	}
	return nil
}

func GetProjectSkills(obj *model.Project) ([]*model.Skill, error) {
	skills := []*model.Skill{}
	allSkills, status := GetSkills()
	if status != 200 {
		return nil, errors.New("failed to get all skills")
	}
	for _, skillId := range obj.SkillIds {
		found, index := utils.SkillBinarySearch(allSkills, &model.Skill{ID: skillId})
		if !found {
			return nil, errors.New("failed to find skill with id: " + skillId + " in project: " + obj.Name)
		}
		skills = append(skills, allSkills[index])
	}
	return skills, nil
}
