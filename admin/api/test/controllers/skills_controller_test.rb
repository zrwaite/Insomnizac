require "test_helper"

class SkillsControllerTest < ActionDispatch::IntegrationTest
  setup do
    @skill = skills(:one)
  end

  test "should get index" do
    get skills_url, as: :json
    assert_response :success
  end

  test "should create skill" do
    assert_difference("Skill.count") do
      post skills_url, params: { skill: {  } }, as: :json
    end

    assert_response :created
  end

  test "should show skill" do
    get skill_url(@skill), as: :json
    assert_response :success
  end

  test "should update skill" do
    patch skill_url(@skill), params: { skill: {  } }, as: :json
    assert_response :success
  end

  test "should destroy skill" do
    assert_difference("Skill.count", -1) do
      delete skill_url(@skill), as: :json
    end

    assert_response :no_content
  end
end
