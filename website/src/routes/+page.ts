import { error } from '@sveltejs/kit';
import type { ProjectType } from 'src/types';
import { graphql } from '$lib/data/graphql';
import { PROJECTS_QUERY } from '$lib/queries/projects';

export interface HomeData {	
	projects: ProjectType[];
}

export async function load(): Promise<HomeData> {
	try {
		const data = await graphql.request(PROJECTS_QUERY);
		return data;
	} catch (e) {
		throw error(400, 'Request failed');
	}
}