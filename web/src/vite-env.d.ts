
declare module "*.vue";
interface Window {
    remount: any;
    unmount: any;
  readonly '__MICRO_APP_BASE_ROUTE__': string;
  __MICRO_APP_ENVIRONMENT__: any;
}

interface ImportMeta {
  env: {
    BASE_URL: string;
    // Add other environment variables here
  };
}
