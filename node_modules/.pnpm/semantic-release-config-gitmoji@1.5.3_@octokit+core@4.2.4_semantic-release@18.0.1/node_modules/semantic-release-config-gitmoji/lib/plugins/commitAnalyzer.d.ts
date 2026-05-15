import type { ReleaseRule } from '../type';
import type { PluginSpec } from 'semantic-release';
/**
 * commit analyzer
 * @param releaseRules
 */
declare const commitAnalyzer: (releaseRules?: ReleaseRule[]) => PluginSpec;
export default commitAnalyzer;
