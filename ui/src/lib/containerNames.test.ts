import { describe, it, expect } from 'vitest';
import { parseNumberInput, buildContainerNames } from './containerNames';

describe('parseNumberInput', () => {
    it('should parse a single number', () => {
        expect(parseNumberInput('1')).toEqual([1]);
        expect(parseNumberInput('42')).toEqual([42]);
    });

    it('should parse comma-separated numbers', () => {
        expect(parseNumberInput('1,2,3')).toEqual([1, 2, 3]);
        expect(parseNumberInput('1, 2, 3')).toEqual([1, 2, 3]);
    });

    it('should parse a range', () => {
        expect(parseNumberInput('1-5')).toEqual([1, 2, 3, 4, 5]);
        expect(parseNumberInput('3-5')).toEqual([3, 4, 5]);
    });

    it('should parse mixed numbers and ranges', () => {
        expect(parseNumberInput('1,3-5,7')).toEqual([1, 3, 4, 5, 7]);
        expect(parseNumberInput('1, 3-5, 7')).toEqual([1, 3, 4, 5, 7]);
    });

    it('should handle multiple ranges', () => {
        expect(parseNumberInput('1-3,5-7')).toEqual([1, 2, 3, 5, 6, 7]);
    });

    it('should handle whitespace', () => {
        expect(parseNumberInput(' 1 , 2 , 3 ')).toEqual([1, 2, 3]);
        expect(parseNumberInput('1 - 3')).toEqual([1, 2, 3]);
    });

    it('should handle a single-element range', () => {
        expect(parseNumberInput('5-5')).toEqual([5]);
    });

    it('should return empty array for empty input', () => {
        expect(parseNumberInput('')).toEqual([]);
        expect(parseNumberInput('   ')).toEqual([]);
    });

    it('should ignore invalid entries', () => {
        expect(parseNumberInput('1,abc,3')).toEqual([1, 3]);
        expect(parseNumberInput('5-3')).toEqual([]); // invalid range (start > end) is skipped
    });

    it('should handle complex mixed input', () => {
        expect(parseNumberInput('1,2,5-8,10,12-14')).toEqual([1, 2, 5, 6, 7, 8, 10, 12, 13, 14]);
    });
});

describe('buildContainerNames', () => {
    describe('containers type', () => {
        it('should return string numbers', () => {
            expect(buildContainerNames('1,2,3', 'containers')).toEqual(['1', '2', '3']);
        });

        it('should handle ranges', () => {
            expect(buildContainerNames('1,3-5,7', 'containers')).toEqual(['1', '3', '4', '5', '7']);
        });
    });

    describe('bowls type', () => {
        it('should prefix with bowl_', () => {
            expect(buildContainerNames('1,2,3', 'bowls')).toEqual(['bowl_1', 'bowl_2', 'bowl_3']);
        });

        it('should handle ranges', () => {
            expect(buildContainerNames('1,3-5', 'bowls')).toEqual(['bowl_1', 'bowl_3', 'bowl_4', 'bowl_5']);
        });
    });

    describe('jars type', () => {
        it('should prefix with jar_', () => {
            expect(buildContainerNames('1,2,3', 'jars')).toEqual(['jar_1', 'jar_2', 'jar_3']);
        });

        it('should handle ranges', () => {
            expect(buildContainerNames('1,3-5', 'jars')).toEqual(['jar_1', 'jar_3', 'jar_4', 'jar_5']);
        });
    });

    it('should return empty array for empty input', () => {
        expect(buildContainerNames('', 'containers')).toEqual([]);
        expect(buildContainerNames('', 'bowls')).toEqual([]);
        expect(buildContainerNames('', 'jars')).toEqual([]);
    });
});
