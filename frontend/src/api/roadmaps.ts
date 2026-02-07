import client from './client';

export type RoadmapCategory = 'language' | 'framework' | 'skill' | 'project' | 'other';
export type RoadmapStatus = 'active' | 'completed';

export interface RoadmapStep {
  id: number;
  roadmap_id: number;
  title: string;
  description: string;
  order_index: number;
  is_completed: boolean;
  completed_at: string | null;
  resource_url: string;
  created_at: string;
  updated_at: string;
}

export interface Roadmap {
  id: number;
  user_id: number;
  user?: {
    id: number;
    name: string;
    avatar_url: string;
  };
  title: string;
  description: string;
  category: RoadmapCategory;
  is_public: boolean;
  step_count: number;
  completed_step_count: number;
  progress: number;
  status: RoadmapStatus;
  created_at: string;
  updated_at: string;
  completed_at: string | null;
  steps?: RoadmapStep[];
}

export interface RoadmapStats {
  total_roadmaps: number;
  active_roadmaps: number;
  completed_roadmaps: number;
  total_steps: number;
  completed_steps: number;
}

export interface CreateRoadmapRequest {
  title: string;
  description?: string;
  category?: RoadmapCategory;
  is_public?: boolean;
}

export interface UpdateRoadmapRequest {
  title?: string;
  description?: string;
  category?: RoadmapCategory;
  is_public?: boolean;
  status?: RoadmapStatus;
}

export interface CreateStepRequest {
  title: string;
  description?: string;
  resource_url?: string;
  order_index?: number;
}

export interface UpdateStepRequest {
  title?: string;
  description?: string;
  resource_url?: string;
  is_completed?: boolean;
}

export interface ReorderStepsRequest {
  orders: Array<{
    step_id: number;
    order_index: number;
  }>;
}

// Roadmap APIs
export const createRoadmap = (data: CreateRoadmapRequest) =>
  client.post<Roadmap>('/roadmaps', data);

export const updateRoadmap = (id: number, data: UpdateRoadmapRequest) =>
  client.put<Roadmap>(`/roadmaps/${id}`, data);

export const deleteRoadmap = (id: number) =>
  client.delete(`/roadmaps/${id}`);

export const getRoadmap = (id: number) =>
  client.get<Roadmap>(`/roadmaps/${id}`);

export const getMyRoadmaps = () =>
  client.get<Roadmap[]>('/roadmaps');

export const getPublicRoadmaps = (limit = 20, offset = 0) =>
  client.get<{ roadmaps: Roadmap[]; total: number }>('/roadmaps/public', {
    params: { limit, offset },
  });

export const copyRoadmap = (id: number) =>
  client.post<Roadmap>(`/roadmaps/${id}/copy`);

// Step APIs
export const createStep = (roadmapId: number, data: CreateStepRequest) =>
  client.post<RoadmapStep>(`/roadmaps/${roadmapId}/steps`, data);

export const updateStep = (roadmapId: number, stepId: number, data: UpdateStepRequest) =>
  client.put<RoadmapStep>(`/roadmaps/${roadmapId}/steps/${stepId}`, data);

export const deleteStep = (roadmapId: number, stepId: number) =>
  client.delete(`/roadmaps/${roadmapId}/steps/${stepId}`);

export const reorderSteps = (roadmapId: number, data: ReorderStepsRequest) =>
  client.put(`/roadmaps/${roadmapId}/steps/reorder`, data);
