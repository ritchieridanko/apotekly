interface APIResponse<T = undefined> {
  status: number;
  message: string;
  data?: T;
}
