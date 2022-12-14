use yew::{function_component, Html, Properties, html, Callback, MouseEvent};
use yew_router::prelude::use_navigator;

use crate::{models::Project, pages::Route};

#[derive(PartialEq, Properties)]
pub struct ProjectPanelProps {
    pub project: Project
}


#[function_component(ProjectPanel)]
pub fn project_panel(props: &ProjectPanelProps) -> Html {
    let ProjectPanelProps { project } = props;
	let navigator = use_navigator().unwrap();
	let slug = project.slug.clone();

	let edit_button: Callback<MouseEvent> = {
		let slug = slug.clone();
		Callback::from(move |_| {
			navigator.push(&Route::EditProject { slug: slug.clone() });
		})
    };

	html! {
		<div class="project">
			<div class="header">
				<p></p>
				<h3>{project.name.to_owned()}</h3>
				<button onclick={edit_button}>{"Edit"}</button>
			</div>
			<div class="image"><img src={project.image.to_owned()}/></div>
			<a href={project.devpost_link.to_owned()}>{project.devpost_link.to_owned()}</a>
			<a href={project.project_link.to_owned()}>{project.project_link.to_owned()}</a>
			<div>{"Featured: "}{project.featured.to_owned()}</div>
			<div>
				<h4>{"Skills: "}</h4>
				<div class="skills-grid">
					{for project.skills.iter().map(|skill| {
						html! {
							<img src={skill.image.to_owned()}/>
						}})
					}
				</div>
			</div>
			<div>{"Created At: "}{project.created_at.to_owned()}</div>
			<div>{"Updated At: "}{project.updated_at.to_owned()}</div>

		</div>
	}
}